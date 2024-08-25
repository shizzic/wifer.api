package main

import "wifer/server/routes"

func main() {
	// test()
	routes.Declare(&props)
	run()
}

// func test() {
// 	opts := options.New(
// 		"mongorestore",
// 		"100.10",
// 		"",
// 		mongorestore.Usage,
// 		true,
// 		options.EnabledOptions{Auth: true, Connection: true, Namespace: true},
// 	)

// 	inputOpts := &mongorestore.InputOptions{}
// 	opts.AddOptions(inputOpts)
// 	outputOpts := &mongorestore.OutputOptions{}
// 	opts.AddOptions(outputOpts)
// 	targetDir := util.ToUniversalPath("C:\\Users\\punch\\OneDrive\\Рабочий стол\\init_dump")

// 	// connect directly, unless a replica set name is explicitly specified
// 	_, setName := util.ParseConnectionString(opts.Host)
// 	opts.Direct = (setName == "")
// 	opts.ReplicaSetName = setName

// 	provider, err := db.NewSessionProvider(*opts)
// 	if err != nil {
// 		log.Logf(log.Always, "error connecting to host: %v", err)
// 		os.Exit(util.ExitError)
// 	}
// 	provider.SetBypassDocumentValidation(outputOpts.BypassDocumentValidation)

// 	// disable TCP timeouts for restore jobs
// 	provider.SetFlags(db.DisableSocketTimeout)
// 	restore := mongorestore.MongoRestore{
// 		ToolOptions:     opts,
// 		OutputOptions:   outputOpts,
// 		InputOptions:    inputOpts,
// 		TargetDirectory: targetDir,
// 		SessionProvider: provider,
// 	}

// 	if err = restore.Restore(); err != nil {
// 		log.Logf(log.Always, "Failed: %v", err)
// 		if err == util.ErrTerminated {
// 			os.Exit(util.ExitKill)
// 		}
// 		os.Exit(util.ExitError)
// 	}
// }
