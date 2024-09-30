
# 代码规范检查
.PHONY: lint
lint:
	@docker run --rm  \
	  --env PROJECT_TYPE="make-lint" \
	  -v $(shell pwd):/opt/app -w /opt/app \
	  registry.cn-hangzhou.aliyuncs.com/mudu/all-in-one:v1.0.0


# 单元测试
.PHONY: test
test:
	go test -v ./... -cover


# 代码格式化
.PHONY: fmt
fmt:
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
		goimports -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret


# 多环境支持，合并开发分支代码到t[1,2,3]，遇到冲突请手动解决冲突
t%:
	@echo "当前对应环境名称: env0$* 分支：t$*";
	- git branch -D t$*;
	git fetch;
	export branch=`git branch | grep \* | grep -Eo ' .+'` && \
		echo "当前分支: $$branch" && \
		git checkout t$* && \
		git pull --rebase && \
		git merge origin/master && \
		echo "merge: \033[0;31morigin/master\033[0m" && \
		git merge $$branch && \
		echo "merge: \033[0;31m$$branch\033[0m" && \
		git push && \
		git checkout $$branch;

rebase:
	export branch=`git branch | grep \* | grep -Eo ' .+'` && \
		git checkout master && \
		git pull --rebase && \
		git checkout $$branch && \
		git rebase master;


# 放弃本地修改
drop:
	git add .; git stash; git stash drop
