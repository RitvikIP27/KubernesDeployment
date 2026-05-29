# I am building my CI WORKFLOW into .github.workflows.ci.yml
 ##Structure is as follows
 Name of workflow
            jobs (of workflow)
               build: job 1
                   uses-actions/checout@v4(means actions repo se features use jo code ko cehckout krke run kre)
                   (basically ek ubuntu runner me code clone)
                   uses: actions setupx (for dockerizing)
NOTE: this structure each action trigger is makinga runner to be run for  it like for for docker we used docker/setup buildexaction v4 is ALSO  an action 