package app

import (
	"github.com/jsli/ota/radio/app/ota_job"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(func() {
		createion_job := ota_job.ReleaseCreationJob{}
		jobs.Schedule("@every 15s", &createion_job)

		remove_job := ota_job.ReleaseRemoveJob{}
		jobs.Schedule("@every 60s", &remove_job)
	})
}
