package cibt

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var xmlBody = `
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:crif-message:2006-08-23" xmlns:urn1="urn:crif-messagegateway:2006-08-23">
    <soap:Header>
        <urn:Message GId="" MId="" MTs="2020-07-04T15:36:56">
            <!--Optional:-->
            <urn:C UD="?" UId="JUDI0176" UPwd="ss7AYBIqI0">
                <!--You may enter ANY elements at this point-->
            </urn:C>
            
            <!--Zero or more repetitions:-->
            <urn:P SId="CB" PId="SD_Req" PNs="urn:SD_Req.2012-04-27.000">
                <!--You may enter ANY elements at this point-->
         
            </urn:P>
            <!--Optional:-->
            <urn:Tx TxNs="urn:crif-messagegateway:2006-08-23">
                <!--You may enter ANY elements at this point-->
         
            </urn:Tx>
            <!--You may enter ANY elements at this point-->
      
        </urn:Message>
    </soap:Header>
    <soap:Body>
        <urn1:MGRequest>
            <![CDATA[
                  <SD_Req_Input><BasicSearch><SubjectType>1</SubjectType><INN>%s</INN></BasicSearch>
                            </SD_Req_Input>
         ]]>
      
    </urn1:MGRequest>
</soap:Body>
</soap:Envelope>
`

func GetCibtId(inn string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Post("https://a2a.cibt.tj", "text/xml", strings.NewReader(fmt.Sprintf(xmlBody, inn)))
	if err != nil {
		fmt.Println("Error in request -", err)
		return ""
	}
	defer resp.Body.Close()

	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error in reading body -", err)
			return ""
		}
		bodyString = string(bodyBytes)
	}

	start := strings.Index(bodyString, "&lt;CBSubjectCode&gt;")
	if start == -1 {
		fmt.Println(err)
		return ""
	}
	end := strings.Index(bodyString, "&lt;/CBSubjectCode&gt")
	if end == -1 {
		fmt.Println(err)
		return ""
	}
	cibt_id := bodyString[start+21 : end]

	return fmt.Sprint(cibt_id)
}
