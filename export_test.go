package lambroll

var (
	CreateZipArchive  = createZipArchive
	ExpandExcludeFile = expandExcludeFile
	LoadZipArchive    = loadZipArchive
	MergeTags         = mergeTags
	FillDefaultValues = fillDefaultValues
	JSONStr           = jsonStr
	MarshalJSON       = marshalJSON
	NewFunctionFrom   = newFunctionFrom
)

type VersionsOutput = versionsOutput
type VersionsOutputs = versionsOutputs

func (app *App) CallerIdentity() *CallerIdentity {
	return app.callerIdentity
}

func (app *App) LoadFunction(f string) (*Function, error) {
	return app.loadFunction(f)
}
