// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accountsettings "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/accountsettings"
	accountsubscription "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/accountsubscription"
	analysis "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/analysis"
	custompermissions "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/custompermissions"
	dashboard "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/dashboard"
	dataset "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/dataset"
	datasource "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/datasource"
	folder "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/folder"
	foldermembership "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/foldermembership"
	group "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/group"
	groupmembership "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/groupmembership"
	iampolicyassignment "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/iampolicyassignment"
	ingestion "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/ingestion"
	iprestriction "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/iprestriction"
	keyregistration "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/keyregistration"
	quicksightnamespace "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/quicksightnamespace"
	refreshschedule "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/refreshschedule"
	rolecustompermission "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/rolecustompermission"
	rolemembership "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/rolemembership"
	template "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/template"
	templatealias "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/templatealias"
	theme "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/theme"
	user "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/user"
	usercustompermission "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/usercustompermission"
	vpcconnection "github.com/upbound/provider-aws/v2/internal/controller/namespaced/quicksight/vpcconnection"
)

// Setup_quicksight creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_quicksight(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accountsettings.Setup,
		accountsubscription.Setup,
		analysis.Setup,
		custompermissions.Setup,
		dashboard.Setup,
		dataset.Setup,
		datasource.Setup,
		folder.Setup,
		foldermembership.Setup,
		group.Setup,
		groupmembership.Setup,
		iampolicyassignment.Setup,
		ingestion.Setup,
		iprestriction.Setup,
		keyregistration.Setup,
		quicksightnamespace.Setup,
		refreshschedule.Setup,
		rolecustompermission.Setup,
		rolemembership.Setup,
		template.Setup,
		templatealias.Setup,
		theme.Setup,
		user.Setup,
		usercustompermission.Setup,
		vpcconnection.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_quicksight creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_quicksight(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accountsettings.SetupGated,
		accountsubscription.SetupGated,
		analysis.SetupGated,
		custompermissions.SetupGated,
		dashboard.SetupGated,
		dataset.SetupGated,
		datasource.SetupGated,
		folder.SetupGated,
		foldermembership.SetupGated,
		group.SetupGated,
		groupmembership.SetupGated,
		iampolicyassignment.SetupGated,
		ingestion.SetupGated,
		iprestriction.SetupGated,
		keyregistration.SetupGated,
		quicksightnamespace.SetupGated,
		refreshschedule.SetupGated,
		rolecustompermission.SetupGated,
		rolemembership.SetupGated,
		template.SetupGated,
		templatealias.SetupGated,
		theme.SetupGated,
		user.SetupGated,
		usercustompermission.SetupGated,
		vpcconnection.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
