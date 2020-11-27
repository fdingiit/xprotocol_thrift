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

// DeleteNodeLabel invokes the retailcloud.DeleteNodeLabel API synchronously
// api document: https://help.aliyun.com/api/retailcloud/deletenodelabel.html
func (client *Client) DeleteNodeLabel(request *DeleteNodeLabelRequest) (response *DeleteNodeLabelResponse, err error) {
	response = CreateDeleteNodeLabelResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteNodeLabelWithChan invokes the retailcloud.DeleteNodeLabel API asynchronously
// api document: https://help.aliyun.com/api/retailcloud/deletenodelabel.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteNodeLabelWithChan(request *DeleteNodeLabelRequest) (<-chan *DeleteNodeLabelResponse, <-chan error) {
	responseChan := make(chan *DeleteNodeLabelResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteNodeLabel(request)
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

// DeleteNodeLabelWithCallback invokes the retailcloud.DeleteNodeLabel API asynchronously
// api document: https://help.aliyun.com/api/retailcloud/deletenodelabel.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteNodeLabelWithCallback(request *DeleteNodeLabelRequest, callback func(response *DeleteNodeLabelResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteNodeLabelResponse
		var err error
		defer close(result)
		response, err = client.DeleteNodeLabel(request)
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

// DeleteNodeLabelRequest is the request struct for api DeleteNodeLabel
type DeleteNodeLabelRequest struct {
	*requests.RpcRequest
	LabelKey   string           `position:"Query" name:"LabelKey"`
	LabelValue string           `position:"Query" name:"LabelValue"`
	Force      requests.Boolean `position:"Query" name:"Force"`
	ClusterId  string           `position:"Query" name:"ClusterId"`
}

// DeleteNodeLabelResponse is the response struct for api DeleteNodeLabel
type DeleteNodeLabelResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	ErrMsg    string `json:"ErrMsg" xml:"ErrMsg"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
}

// CreateDeleteNodeLabelRequest creates a request to invoke DeleteNodeLabel API
func CreateDeleteNodeLabelRequest() (request *DeleteNodeLabelRequest) {
	request = &DeleteNodeLabelRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("retailcloud", "2018-03-13", "DeleteNodeLabel", "", "")
	request.Method = requests.POST
	return
}

// CreateDeleteNodeLabelResponse creates a response to parse from DeleteNodeLabel response
func CreateDeleteNodeLabelResponse() (response *DeleteNodeLabelResponse) {
	response = &DeleteNodeLabelResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
