package codec

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	raw, _ := DecodeBase64("kPH+bIxk5D2deZiIxcaaaA==")

	aced, _ := DecodeHex("aced0005737200176a6176612e7574696c2e5072696f72697479517565756594da30b4fb3f82b103000249000473697a654c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f727870000000027372002b6f72672e6170616368652e636f6d6d6f6e732e6265616e7574696c732e4265616e436f6d70617261746f72e3a188ea7322a4480200024c000870726f70657274797400114c6a6176612f6c616e672f537472696e674c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f7278707400106f757470757450726f706572746965737372003f6f72672e6170616368652e636f6d6d6f6e732e636f6c6c656374696f6e732e636f6d70617261746f72732e436f6d70617261626c65436f6d70617261746f72fbf49925b86eb13702000078707704000000037372003a636f6d2e73756e2e6f72672e6170616368652e78616c616e2e696e7465726e616c2e78736c74632e747261782e54656d706c61746573496d706c09574fc16eacab3303000649000d5f696e64656e744e756d62657249000e5f7472616e736c6574496e6465785b000a5f62797465636f6465737400035b5b425b00065f636c6173737400125b4c6a6176612f6c616e672f436c6173733b4c00055f6e616d657400124c6a6176612f6c616e672f537472696e673b4c00115f6f757470757450726f706572746965737400164c6a6176612f7574696c2f50726f706572746965733b787000000000ffffffff757200035b5b424bfd19156767db37020000787000000001757200025b42acf317f8060854e002000078700000066acafebabe0000003400390a0003002207003707002507002601001073657269616c56657273696f6e5549440100014a01000d436f6e7374616e7456616c756505ad2093f391ddef3e0100063c696e69743e010003282956010004436f646501000f4c696e654e756d6265725461626c650100124c6f63616c5661726961626c655461626c65010004746869730100105472616e736c6174655061796c6f616401000c496e6e6572436c617373657301001e4c436c6173734c6f61646572245472616e736c6174655061796c6f61643b0100097472616e73666f726d010072284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b5b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b2956010008646f63756d656e7401002d4c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b01000868616e646c6572730100425b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a457863657074696f6e730700270100a6284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b29560100086974657261746f720100354c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b01000768616e646c65720100414c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a536f7572636546696c65010010436c6173734c6f616465722e6a6176610c000a000b07002801001c436c6173734c6f61646572245472616e736c6174655061796c6f6164010040636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f72756e74696d652f41627374726163745472616e736c65740100146a6176612f696f2f53657269616c697a61626c65010039636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f5472616e736c6574457863657074696f6e01000b436c6173734c6f616465720100083c636c696e69743e0100116a6176612f6c616e672f52756e74696d6507002a01000a67657452756e74696d6501001528294c6a6176612f6c616e672f52756e74696d653b0c002c002d0a002b002e01004962617368202d63207b6563686f2c64473931593267674c335274634339695a5746756458527062484d784c574a686332673d7d7c7b6261736536342c2d647d7c7b626173682c2d697d08003001000465786563010027284c6a6176612f6c616e672f537472696e673b294c6a6176612f6c616e672f50726f636573733b0c003200330a002b003401000d537461636b4d61705461626c65010004746573740100064c746573743b002100020003000100040001001a000500060001000700000002000800040001000a000b0001000c0000002f00010001000000052ab70001b100000002000d00000006000100000015000e0000000c000100000005000f003800000001001300140002000c0000003f0000000300000001b100000002000d0000000600010000001b000e00000020000300000001000f0038000000000001001500160001000000010017001800020019000000040001001a00010013001b0002000c000000490000000400000001b100000002000d00000006000100000020000e0000002a000400000001000f003800000000000100150016000100000001001c001d000200000001001e001f00030019000000040001001a00080029000b0001000c00000024000300020000000fa70003014cb8002f1231b6003557b1000000010036000000030001030002002000000002002100110000000a00010002002300100009707400144869567671664a79516f45647553745442485861707701007871007e000f78")
	data, err := AESCBCEncrypt(aced, raw, nil)
	if err != nil {
		return
	}
	_ = data
	//println(strconv.Quote(string(data)))
	//println(strconv.Quote(string(iv)))
	//
	//println(EncodeBase64(append(iv, data...)))
}

func TestAESGCMEnc(t *testing.T) {
	raw, _ := DecodeBase64("kPH+bIxk5D2deZiIxcaaaA==")
	data, err := AESGCMEncrypt(raw, "asdfasdfasdfasdfasdfaaa", nil)
	if err != nil {
		panic(err)
	}
	spew.Dump(data)
	originData, err := AESGCMDecrypt(raw, data, nil)
	if err != nil {
		panic(err)
	}
	spew.Dump(originData)
}

//func TestAESGCMEnc1(t *testing.T) {
//	raw, _ := DecodeBase64("kPH+bIxk5D2deZiIxcaaaA==")
//	originData, err := AESGCMDec(raw, []byte("as23423423asdfasdfasdfasdfasdfasdfas"), nil)
//	if err != nil {
//		panic(err)
//	}
//	spew.Dump(originData)
//}

func TestAesCBC(t *testing.T) {
	secret := []byte("1231231231231234")
	tdata := "aced0005737200176a6176612e7574696c2e5072696f72697479517565756594da30b4fb3f82b103000249000473697a654c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f727870000000027372002b6f72672e6170616368652e636f6d6d6f6e732e6265616e7574696c732e4265616e436f6d70617261746f72e3a188ea7322a4480200024c000870726f70657274797400114c6a6176612f6c616e672f537472696e674c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f7278707400106f757470757450726f706572746965737372003f6f72672e6170616368652e636f6d6d6f6e732e636f6c6c656374696f6e732e636f6d70617261746f72732e436f6d70617261626c65436f6d70617261746f72fbf49925b86eb13702000078707704000000037372003a636f6d2e73756e2e6f72672e6170616368652e78616c616e2e696e7465726e616c2e78736c74632e747261782e54656d706c61746573496d706c09574fc16eacab3303000649000d5f696e64656e744e756d62657249000e5f7472616e736c6574496e6465785b000a5f62797465636f6465737400035b5b425b00065f636c6173737400125b4c6a6176612f6c616e672f436c6173733b4c00055f6e616d657400124c6a6176612f6c616e672f537472696e673b4c00115f6f757470757450726f706572746965737400164c6a6176612f7574696c2f50726f706572746965733b787000000000ffffffff757200035b5b424bfd19156767db37020000787000000001757200025b42acf317f8060854e002000078700000066acafebabe0000003400390a0003002207003707002507002601001073657269616c56657273696f6e5549440100014a01000d436f6e7374616e7456616c756505ad2093f391ddef3e0100063c696e69743e010003282956010004436f646501000f4c696e654e756d6265725461626c650100124c6f63616c5661726961626c655461626c65010004746869730100105472616e736c6174655061796c6f616401000c496e6e6572436c617373657301001e4c436c6173734c6f61646572245472616e736c6174655061796c6f61643b0100097472616e73666f726d010072284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b5b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b2956010008646f63756d656e7401002d4c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b01000868616e646c6572730100425b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a457863657074696f6e730700270100a6284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b29560100086974657261746f720100354c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b01000768616e646c65720100414c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a536f7572636546696c65010010436c6173734c6f616465722e6a6176610c000a000b07002801001c436c6173734c6f61646572245472616e736c6174655061796c6f6164010040636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f72756e74696d652f41627374726163745472616e736c65740100146a6176612f696f2f53657269616c697a61626c65010039636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f5472616e736c6574457863657074696f6e01000b436c6173734c6f616465720100083c636c696e69743e0100116a6176612f6c616e672f52756e74696d6507002a01000a67657452756e74696d6501001528294c6a6176612f6c616e672f52756e74696d653b0c002c002d0a002b002e01004962617368202d63207b6563686f2c64473931593267674c335274634339695a5746756458527062484d784c574a686332673d7d7c7b6261736536342c2d647d7c7b626173682c2d697d08003001000465786563010027284c6a6176612f6c616e672f537472696e673b294c6a6176612f6c616e672f50726f636573733b0c003200330a002b003401000d537461636b4d61705461626c65010004746573740100064c746573743b002100020003000100040001001a000500060001000700000002000800040001000a000b0001000c0000002f00010001000000052ab70001b100000002000d00000006000100000015000e0000000c000100000005000f003800000001001300140002000c0000003f0000000300000001b100000002000d0000000600010000001b000e00000020000300000001000f0038000000000001001500160001000000010017001800020019000000040001001a00010013001b0002000c000000490000000400000001b100000002000d00000006000100000020000e0000002a000400000001000f003800000000000100150016000100000001001c001d000200000001001e001f00030019000000040001001a00080029000b0001000c00000024000300020000000fa70003014cb8002f1231b6003557b1000000010036000000030001030002002000000002002100110000000a00010002002300100009707400144869567671664a79516f45647553745442485861707701007871007e000f78"
	res, err := AESCBCEncrypt(secret, interfaceToBytes(tdata), nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	result, err := AESCBCDecrypt(secret, res, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.True(t, string(result) == tdata)
}

func TestAesECB(t *testing.T) {
	secret := []byte("1231231231231234")
	tdata := "aced0005737200176a6176612e7574696c2e5072696f72697479517565756594da30b4fb3f82b103000249000473697a654c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f727870000000027372002b6f72672e6170616368652e636f6d6d6f6e732e6265616e7574696c732e4265616e436f6d70617261746f72e3a188ea7322a4480200024c000870726f70657274797400114c6a6176612f6c616e672f537472696e674c000a636f6d70617261746f727400154c6a6176612f7574696c2f436f6d70617261746f7278707400106f757470757450726f706572746965737372003f6f72672e6170616368652e636f6d6d6f6e732e636f6c6c656374696f6e732e636f6d70617261746f72732e436f6d70617261626c65436f6d70617261746f72fbf49925b86eb13702000078707704000000037372003a636f6d2e73756e2e6f72672e6170616368652e78616c616e2e696e7465726e616c2e78736c74632e747261782e54656d706c61746573496d706c09574fc16eacab3303000649000d5f696e64656e744e756d62657249000e5f7472616e736c6574496e6465785b000a5f62797465636f6465737400035b5b425b00065f636c6173737400125b4c6a6176612f6c616e672f436c6173733b4c00055f6e616d657400124c6a6176612f6c616e672f537472696e673b4c00115f6f757470757450726f706572746965737400164c6a6176612f7574696c2f50726f706572746965733b787000000000ffffffff757200035b5b424bfd19156767db37020000787000000001757200025b42acf317f8060854e002000078700000066acafebabe0000003400390a0003002207003707002507002601001073657269616c56657273696f6e5549440100014a01000d436f6e7374616e7456616c756505ad2093f391ddef3e0100063c696e69743e010003282956010004436f646501000f4c696e654e756d6265725461626c650100124c6f63616c5661726961626c655461626c65010004746869730100105472616e736c6174655061796c6f616401000c496e6e6572436c617373657301001e4c436c6173734c6f61646572245472616e736c6174655061796c6f61643b0100097472616e73666f726d010072284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b5b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b2956010008646f63756d656e7401002d4c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b01000868616e646c6572730100425b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a457863657074696f6e730700270100a6284c636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f444f4d3b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b4c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b29560100086974657261746f720100354c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f64746d2f44544d417869734974657261746f723b01000768616e646c65720100414c636f6d2f73756e2f6f72672f6170616368652f786d6c2f696e7465726e616c2f73657269616c697a65722f53657269616c697a6174696f6e48616e646c65723b01000a536f7572636546696c65010010436c6173734c6f616465722e6a6176610c000a000b07002801001c436c6173734c6f61646572245472616e736c6174655061796c6f6164010040636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f72756e74696d652f41627374726163745472616e736c65740100146a6176612f696f2f53657269616c697a61626c65010039636f6d2f73756e2f6f72672f6170616368652f78616c616e2f696e7465726e616c2f78736c74632f5472616e736c6574457863657074696f6e01000b436c6173734c6f616465720100083c636c696e69743e0100116a6176612f6c616e672f52756e74696d6507002a01000a67657452756e74696d6501001528294c6a6176612f6c616e672f52756e74696d653b0c002c002d0a002b002e01004962617368202d63207b6563686f2c64473931593267674c335274634339695a5746756458527062484d784c574a686332673d7d7c7b6261736536342c2d647d7c7b626173682c2d697d08003001000465786563010027284c6a6176612f6c616e672f537472696e673b294c6a6176612f6c616e672f50726f636573733b0c003200330a002b003401000d537461636b4d61705461626c65010004746573740100064c746573743b002100020003000100040001001a000500060001000700000002000800040001000a000b0001000c0000002f00010001000000052ab70001b100000002000d00000006000100000015000e0000000c000100000005000f003800000001001300140002000c0000003f0000000300000001b100000002000d0000000600010000001b000e00000020000300000001000f0038000000000001001500160001000000010017001800020019000000040001001a00010013001b0002000c000000490000000400000001b100000002000d00000006000100000020000e0000002a000400000001000f003800000000000100150016000100000001001c001d000200000001001e001f00030019000000040001001a00080029000b0001000c00000024000300020000000fa70003014cb8002f1231b6003557b1000000010036000000030001030002002000000002002100110000000a00010002002300100009707400144869567671664a79516f45647553745442485861707701007871007e000f78"
	res, err := AESECBEncrypt(secret, interfaceToBytes(tdata), nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	result, err := AESECBDecrypt(secret, res, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.True(t, string(result) == tdata)
}
