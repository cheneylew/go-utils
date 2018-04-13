package stock

import (
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2"
	"fmt"
	"log"
	"strings"
)

func pushNotification (msg string)  {
	cert, err := certificate.FromP12File("/Users/dejunliu/Desktop/apn-push-ios1234.p12", "1234")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notification := &apns2.Notification{}
	uuid := strings.Replace("3a5dfd69 2ba6c42d 26358ed4 cd497f73 55a1cd4e 4a1b7f66 9bbf7867 17f3aa1d"," ","",-1)
	notification.DeviceToken = uuid
	notification.Topic = "com.hyh.pushdemo"
	notification.Payload = []byte(fmt.Sprintf(`{"aps":{"alert":"%s"}}`, msg)) // See Payload section below

	client := apns2.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}


	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}