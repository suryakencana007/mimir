/*  aws.go
*
* @Author:             Nanang Suryadi
* @Date:               February 20, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-20 23:00
 */

package request

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
)

func HandleEvent(handler *chi.Mux) interface{} {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		chiLambda := chiadapter.New(handler)
		return chiLambda.Proxy(req)
	}
}
