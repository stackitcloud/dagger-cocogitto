set dotenv-load
set dotenv-filename := ".env"

commit:
    dagger call commit \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN
    
check:
    dagger call check \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

install-hooks:
    dagger call install-hooks \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

verify:
    dagger call verify \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

log:
    dagger call log \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

changelog:
    dagger call changelog \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

bump:
    dagger call bump \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

get-version:
    dagger call get-version \
        --repository-url=$REPOSITORY_URL \
        --user=$USER \
        --git-token=env:GIT_TOKEN

