package gg

import (
	"log"
	"path"
	"time"

	"ulambda/fslib"
	"ulambda/proc"
	"ulambda/procclnt"
)

const (
	THUNK_OUTPUTS_SUFFIX      = ".thunk-outputs"
	EXIT_DEPENDENCIES_SUFFIX  = ".exit-dependencies"
	INPUT_DEPENDENCIES_SUFFIX = ".input-dependencies"
	SHEBANG_DIRECTIVE         = "#!/usr/bin/env gg-force-and-run"
)

type ExecutorLauncher interface {
	FsLambda
	Spawn(*proc.Proc) error
	Started() error
}

type Orchestrator struct {
	pid          proc.Tpid
	cwd          string
	targets      []string
	targetHashes []string
	*fslib.FsLib
	*procclnt.ProcClnt
}

func MakeOrchestrator(args []string, debug bool) (*Orchestrator, error) {
	log.Printf("Orchestrator: %v\n", args)
	orc := &Orchestrator{}

	orc.pid = proc.Tpid(args[0])
	orc.cwd = args[1]
	orc.targets = args[2:]
	fls := fslib.MakeFsLib("orchestrator")
	orc.FsLib = fls
	orc.ProcClnt = procclnt.MakeProcClnt(fls)
	orc.Started()
	return orc, nil
}

func (orc *Orchestrator) Exit() {
	orc.Exited(nil)
}

func (orc *Orchestrator) Work() {
	setUpRemoteDirs(orc)
	copyRemoteDirTree(orc, path.Join(orc.cwd, ".gg"), ggRemote("", ""))
	reductionWriters := []proc.Tpid{}
	start := time.Now()
	for i, target := range orc.targets {
		targetHash := getTargetHash(orc, orc.cwd, target)
		orc.targetHashes = append(orc.targetHashes, targetHash)
		g := orc.ingestStaticGraph(targetHash)
		// Ignore reductions, which aren't actually executable
		if !isReduction(targetHash) {
			orc.executeStaticGraph(targetHash, g)
			rwPid := spawnReductionWriter(orc, orc.targets[i], targetHash, path.Join(orc.cwd, "results"), "", []string{})
			reductionWriters = append(reductionWriters, rwPid)
		}
	}
	for _, rw := range reductionWriters {
		orc.WaitExit(rw)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Elapsed time: %v ms", elapsed.Milliseconds())
}

func (orc *Orchestrator) ingestStaticGraph(targetHash string) *Graph {
	g := MakeGraph()
	var current string
	queue := []string{targetHash}
	// Will loop inifinitely if the graph is circular
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]
		if isReduction(current) {
			current = thunkHashFromReduction(current)
		}
		exitDeps := getExitDependencies(orc, current)
		exitDeps = thunkHashesFromReductions(exitDeps)
		// XXX Should I add the thunk's hash as an output file?
		g.AddThunk(current, exitDeps, []string{})
		queue = append(queue, exitDeps...)
	}
	return g
}

func (orc *Orchestrator) executeStaticGraph(targetHash string, g *Graph) {
	thunks := g.GetThunks()
	for _, thunk := range thunks {
		exitDeps := outputHandlerPids(thunk.deps)
		if reductionExists(orc, thunk.hash) || currentlyExecuting(orc, thunk.hash) || isReduction(thunk.hash) {
			continue
		}
		exPid, err := spawnExecutor(orc, thunk.hash, exitDeps)
		if err != nil {
			log.Fatalf("Error orchestrator: %v", err)
		}
		outputHandlerPid := spawnThunkOutputHandler(orc, []string{exPid.String()}, thunk.hash, []string{thunk.hash})
		spawnNoOp(orc, outputHandlerPid)
	}
}

func (orc *Orchestrator) waitPids(pids []proc.Tpid) {
	for _, p := range pids {
		orc.WaitExit(p)
	}
}

func (orc *Orchestrator) Name() string {
	return "Orchestrator"
}
