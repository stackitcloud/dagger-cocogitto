set dotenv-load
set dotenv-filename := ".env"
    
bump *ARGS:
    dagger call bump \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

changelog *ARGS:
    dagger call changelog \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

check *ARGS:
    dagger call check \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

commit *ARGS:
    dagger call commit \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

get-version *ARGS:
    dagger call get-version \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

install-hooks *ARGS:
    dagger call install-hooks \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

list:
    dagger functions

log *ARGS:
    dagger call log \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}

verify *ARGS:
    dagger call verify \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN {{ ARGS }}