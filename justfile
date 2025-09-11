set dotenv-load
set dotenv-filename := ".env"
    
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
