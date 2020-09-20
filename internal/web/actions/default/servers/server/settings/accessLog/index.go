package accessLog

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("accessLog")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 获取配置
	webResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webResp.Config, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["accessLogConfig"] = webConfig.AccessLog

	// 可选的缓存策略
	policiesResp, err := this.RPC().HTTPAccessLogPolicyRPC().FindAllEnabledHTTPAccessLogPolicies(this.AdminContext(), &pb.FindAllEnabledHTTPAccessLogPoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyMaps := []maps.Map{}
	for _, policy := range policiesResp.AccessLogPolicies {
		policyMaps = append(policyMaps, maps.Map{
			"id":   policy.Id,
			"name": policy.Name,
			"isOn": policy.IsOn, // TODO 这里界面上显示是否开启状态
		})
	}
	this.Data["accessLogPolicies"] = policyMaps

	// 通用变量
	this.Data["fields"] = serverconfigs.HTTPAccessLogFields
	this.Data["defaultFieldCodes"] = serverconfigs.HTTPAccessLogDefaultFieldsCodes

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId         int64
	AccessLogJSON []byte

	Must *actions.Must
}) {
	// TODO 检查参数

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPAccessLog(this.AdminContext(), &pb.UpdateHTTPAccessLogRequest{
		WebId:         params.WebId,
		AccessLogJSON: params.AccessLogJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}