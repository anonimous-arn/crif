package cibt

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var xmlbody2 = `
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:crif-message:2006-08-23" xmlns:urn1="urn:crif-messagegateway:2006-08-23">
        <soap:Header>
            <urn:Message GId="" MId="" MTs="2018-07-04T15:36:56">
                <!--Optional:-->
                <urn:C UD="?" UId="-----" UPwd="------">
                    <!--You may enter ANY elements at this point-->

                </urn:C>
                <!--Zero or more repetitions:-->
                <urn:P SId="CB" PId="CI_Req" PNs="urn:CI_Req.2012-04-27.000">
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
    <CI_Req_Input>
    <CBSubjectCode>%s</CBSubjectCode>
    <Report>
    <Type>1</Type>
    </Report>
    </CI_Req_Input>
             ]]>

        </urn1:MGRequest>
    </soap:Body>
    </soap:Envelope>
`

func GetCibtInfo(cibt_id string) (bodyString string) {
	client := http.Client{
		Timeout: time.Minute * 1,
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", "https://a2a.cibt.tj", strings.NewReader(fmt.Sprintf(xmlbody2, cibt_id)))
	if err != nil {
		fmt.Println(err)
		return bodyString
	}

	request.Header.Set("User-Agent", "My-Client")
	request.Header.Set("Content-Type", "text/xml")
	request.Header.Set("Accept", "text/xml")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return bodyString
	}
	defer resp.Body.Close()

	// var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return bodyString
		}
		bodyString = string(bodyBytes)
	}
	// fmt.Println("bodyString - ", bodyString)

	return bodyString
}
