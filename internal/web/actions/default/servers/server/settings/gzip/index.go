package gzip

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("gzip")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	webConfigResp, err := this.RPC().ServerRPC().FindAndInitServerWebConfig(this.AdminContext(), &pb.FindAndInitServerWebRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(webConfigResp.Config, webConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webResp, err := this.RPC().HTTPWebRPC().FindEnabledHTTPWeb(this.AdminContext(), &pb.FindEnabledHTTPWebRequest{WebId: webConfig.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	web := webResp.Web

	this.Data["webId"] = webConfig.Id

	gzipId := web.GzipId
	gzipConfig := &serverconfigs.HTTPGzipConfig{
		Id:   0,
		IsOn: true,
	}
	if gzipId > 0 {
		resp, err := this.RPC().HTTPGzipRPC().FindEnabledHTTPGzipConfig(this.AdminContext(), &pb.FindEnabledGzipConfigRequest{GzipId: gzipId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = json.Unmarshal(resp.Config, gzipConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["gzipConfig"] = gzipConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	GzipId    int64
	Level     int
	MinLength string
	MaxLength string

	Must *actions.Must
}) {
	if params.Level < 0 || params.Level > 9 {
		this.Fail("请选择正确的压缩级别")
	}

	minLength := &pb.SizeCapacity{Count: -1}
	if len(params.MinLength) > 0 {
		err := json.Unmarshal([]byte(params.MinLength), minLength)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	maxLength := &pb.SizeCapacity{Count: -1}
	if len(params.MaxLength) > 0 {
		err := json.Unmarshal([]byte(params.MaxLength), maxLength)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}


	if params.GzipId > 0 {
		_, err := this.RPC().HTTPGzipRPC().UpdateHTTPGzip(this.AdminContext(), &pb.UpdateHTTPGzipRequest{
			GzipId:    params.GzipId,
			Level:     int32(params.Level),
			MinLength: minLength,
			MaxLength: maxLength,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		resp, err := this.RPC().HTTPGzipRPC().CreateHTTPGzip(this.AdminContext(), &pb.CreateHTTPGzipRequest{
			Level:     0,
			MinLength: nil,
			MaxLength: nil,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		gzipId := resp.GzipId

		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebGzip(this.AdminContext(), &pb.UpdateHTTPWebGzipRequest{
			WebId:  params.WebId,
			GzipId: gzipId,
		})
		if err != nil {
			this.ErrorPage(err)
		}
	}

	this.Success()
}