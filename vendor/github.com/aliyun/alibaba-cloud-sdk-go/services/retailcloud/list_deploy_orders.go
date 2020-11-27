package retailcloud

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ListDeployOrders invokes the retailcloud.ListDeployOrders API synchronously
// api document: https://help.aliyun.com/api/retailcloud/listdeployorders.html
func (client *Client) ListDeployOrders(request *ListDeployOrdersRequest) (response *ListDeployOrdersResponse, err error) {
	response = CreateListDeployOrdersResponse()
	err = client.DoAction(request, response)
	return
}

// ListDeployOrdersWithChan invokes the retailcloud.ListDeployOrders API asynchronously
// api document: https://help.aliyun.com/api/retailcloud/listdeployorders.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListDeployOrdersWithChan(request *ListDeployOrdersRequest) (<-chan *ListDeployOrdersResponse, <-chan error) {
	responseChan := make(chan *ListDeployOrdersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListDeployOrders(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ListDeployOrdersWithCallback invokes the retailcloud.ListDeployOrders API asynchronously
// api document: https://help.aliyun.com/api/retailcloud/listdeployorders.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListDeployOrdersWithCallback(request *ListDeployOrdersRequest, callback func(response *ListDeployOrdersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListDeployOrdersResponse
		var err error
		defer close(result)
		response, err = client.ListDeployOrders(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ListDeployOrdersRequest is the request struct for api ListDeployOrders
type ListDeployOrdersRequest struct {
	*requests.RpcRequest
	StartTimeGreaterThanOrEqualTo string           `position:"Query" name:"StartTimeGreaterThanOrEqualTo"`
	StatusList                    *[]string        `position:"Body" name:"StatusList"  type:"Repeated"`
	EnvId                         requests.Integer `position:"Query" name:"EnvId"`
	EndTimeGreaterThan            string           `position:"Query" name:"EndTimeGreaterThan"`
	PageNumber                    requests.Integer `position:"Query" name:"PageNumber"`
	PauseType                     string           `position:"Query" name:"PauseType"`
	ResultList                    *[]string        `position:"Body" name:"ResultList"  type:"Repeated"`
	StartTimeGreaterThan          string           `position:"Query" name:"StartTimeGreaterThan"`
	StartTimeLessThan             string           `position:"Query" name:"StartTimeLessThan"`
	StartTimeLessThanOrEqualTo    string           `position:"Query" name:"StartTimeLessThanOrEqualTo"`
	AppId                         requests.Integer `position:"Query" name:"AppId"`
	EnvType                       string           `position:"Query" name:"EnvType"`
	PageSize                      requests.Integer `position:"Query" name:"PageSize"`
	EndTimeGreaterThanOrEqualTo   string           `position:"Query" name:"EndTimeGreaterThanOrEqualTo"`
	EndTimeLessThan               string           `position:"Query" name:"EndTimeLessThan"`
	EndTimeLessThanOrEqualTo      string           `position:"Query" name:"EndTimeLessThanOrEqualTo"`
	PartitionType                 string           `position:"Query" name:"PartitionType"`
	DeployCategory                string           `position:"Query" name:"DeployCategory"`
	DeployType                    string           `position:"Query" name:"DeployType"`
	Status                        requests.Integer `position:"Query" name:"Status"`
}

// ListDeployOrdersResponse is the response struct for api ListDeployOrders
type ListDeployOrdersResponse struct {
	*responses.BaseResponse
	Code       int                   `json:"Code" xml:"Code"`
	ErrorMsg   string                `json:"ErrorMsg" xml:"ErrorMsg"`
	PageNumber int                   `json:"PageNumber" xml:"PageNumber"`
	PageSize   int                   `json:"PageSize" xml:"PageSize"`
	RequestId  string                `json:"RequestId" xml:"RequestId"`
	TotalCount int64                 `json:"TotalCount" xml:"TotalCount"`
	Data       []DeployOrderInstance `json:"Data" xml:"Data"`
}

// CreateListDeployOrdersRequest creates a request to invoke ListDeployOrders API
func CreateListDeployOrdersRequest() (request *ListDeployOrdersRequest) {
	request = &ListDeployOrdersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("retailcloud", "2018-03-13", "ListDeployOrders", "", "")
	request.Method = requests.POST
	return
}

// CreateListDeployOrdersResponse creates a response to parse from ListDeployOrders response
func CreateListDeployOrdersResponse() (response *ListDeployOrdersResponse) {
	response = &ListDeployOrdersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
