package decode

import (
	"testing"
)

func TestJson_DecodeReportFromByte(t *testing.T) {
	jsonReport := JsonReport{}
	bs := `{"test2":true,"cid":"47514950895225","scene":"tests","passport":"AYhfmsnCX5rKdpkLdnskuapTN29oVkouFAwAAAAAAAAAAAAAAAAAAAAA.cXR0X2xvZ2luX3dlY2hhdHxhfDMuMTAuMTIuMDAxLjIwMTAyNQ.wGSH1EqVh-RHNz0v6poCwOjHckLC0TDyxg_xMFPF4j8FhAZBYRHB8pX-6oEf39_pNmpwCBLu8RGLnOvgd1D4kw","tk":"ACFYnETqd9I1rmWZWFJ8mWcSexcJmXc7IB40NzUxNDk1MDg5NTIyNQ","tuid":"WJxE6nfSNa5lmVhSfJlnEg","member_id":"","ip":"","app_version":"3.10.12.001.201025","platform":"android","sid":"","status":7,"passport_type":136,"passport_version":1,"valid":1234.5678,"indate":0,"period":180,"issue":"2020-10-29 21:55:14","expire":"2020-10-29 21:58:14","expire_gap":13095882,"ext":{"res":"not_register","tuid":"AAAAAAAAAAAAAAAAAAAAAA"},"timestamp": 1617075776829, "test1": [{"timestamp": 1617075776829, "timestamp2": 16170759}]}`
	//bs := `{"cid":"47514950895225","scene":"tests","passport":"AYhfmsnCX5rKdpkLdnskuapTN29oVkouFAwAAAAAAAAAAAAAAAAAAAAA.cXR0X2xvZ2luX3dlY2hhdHxhfDMuMTAuMTIuMDAxLjIwMTAyNQ.wGSH1EqVh-RHNz0v6poCwOjHckLC0TDyxg_xMFPF4j8FhAZBYRHB8pX-6oEf39_pNmpwCBLu8RGLnOvgd1D4kw","tk":"ACFYnETqd9I1rmWZWFJ8mWcSexcJmXc7IB40NzUxNDk1MDg5NTIyNQ","tuid":"WJxE6nfSNa5lmVhSfJlnEg","member_id":"","ip":"","app_version":"3.10.12.001.201025","platform":"android","sid":"","status":7,"passport_type":136,"passport_version":1,"valid":1234.5678,"indate":0,"period":180,"issue":"2020-10-29 21:55:14","expire":"2020-10-29 21:58:14","expire_gap":13095882,"ext":{"res":"not_register","tuid":"AAAAAAAAAAAAAAAAAAAAAA"},"timestamp": 1617075776829}`
	//bs := `{"cid":"47514950895225","scene":"tests","cid":"47514950895225","scene":"tests","cid":"47514950895225"}`
	m := make(map[string]string, 100)
	_, err := jsonReport.DecodeReportFromByte(m, []byte(bs))
	t.Logf("err:%v", err)
	for k, v := range m {
		t.Logf("%s ===> %s", k, v)
	}
}

func BenchmarkJsonDecode(b *testing.B) {
	jsonReport := JsonReport{}
	bs := `{"cid":"47514950895225","scene":"tests","passport":"AYhfmsnCX5rKdpkLdnskuapTN29oVkouFAwAAAAAAAAAAAAAAAAAAAAA.cXR0X2xvZ2luX3dlY2hhdHxhfDMuMTAuMTIuMDAxLjIwMTAyNQ.wGSH1EqVh-RHNz0v6poCwOjHckLC0TDyxg_xMFPF4j8FhAZBYRHB8pX-6oEf39_pNmpwCBLu8RGLnOvgd1D4kw","tk":"ACFYnETqd9I1rmWZWFJ8mWcSexcJmXc7IB40NzUxNDk1MDg5NTIyNQ","tuid":"WJxE6nfSNa5lmVhSfJlnEg","member_id":"","ip":"","app_version":"3.10.12.001.201025","platform":"android","sid":"","status":7,"passport_type":136,"passport_version":1,"valid":1234.5678,"indate":0,"period":180,"issue":"2020-10-29 21:55:14","expire":"2020-10-29 21:58:14","expire_gap":13095882,"ext":{"res":{"xu":"not_register"},"tuid":"AAAAAAAAAAAAAAAAAAAAAA"},"timestamp": 1617075776829, "test1": [{"timestamp": 1617075776829, "timestamp2": 16170759}]}`
	for i := 0; i < b.N; i++ {
		m := make(map[string]string, 100)
		jsonReport.DecodeReportFromByte(m, []byte(bs))
	}
	//log.Errorf("err:%v", err)
	//for k, v := range m {
	//	log.Infof("%s ===> %s", k, v)
	//}
}
