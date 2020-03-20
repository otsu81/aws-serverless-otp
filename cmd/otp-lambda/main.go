package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pcarrier/gauth/gauth"
)

type Event struct {
	ParameterName string `json:"parameterName"`
}

type OtpResponse struct {
	AccountId     string            `json:"parameterName"`
	OTPChallenges map[string]string `json:"otpChallenges"`
}

func randomFakeNumbersString(n int) string {
	const integerBytes = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = integerBytes[rand.Int63()%int64(len(integerBytes))]
	}
	return string(b)
}

func auth(ctx context.Context, event Event) (OtpResponse, error) {

	sess, err := session.NewSession()

	if err != nil {
		panic(err)
	}
	ssmcl := ssm.New(sess)

	key := event.ParameterName
	decrypt := true
	param, err := ssmcl.GetParameter(&ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: &decrypt,
	})

	var curr string
	var next string
	var prev string

	if err != nil {
		prev = randomFakeNumbersString(6)
		curr = randomFakeNumbersString(6)
		next = randomFakeNumbersString(6)
		log.Printf("Fake numbers generated")
	} else {
		secret := *param.Parameter.Value
		currentTS, _ := gauth.IndexNow()

		var erro error
		prev, curr, next, erro = gauth.Codes(secret, currentTS)
		if erro != nil {
			log.Fatalf("Code: %v", erro)
		}
	}

	return OtpResponse{event.ParameterName, map[string]string{"prev": prev, "curr": curr, "next": next}}, nil
}

func main() {
	lambda.Start(auth)
}
