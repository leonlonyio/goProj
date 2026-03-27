cd /Users/goProj/test/se-take-home-assignment/scripts/

1.Make all scripts executable:

chmod +x build.sh test.sh run.sh demo.sh

2.Build the application:

./build.sh


3.Run tests:

./test.sh


4.Run the simulation (non-interactive mode for GitHub Actions):

./run.sh

5.For interactive mode (for the next round):

./demo.sh

# Fork

1. 访问原仓库：
   https://github.com/feedmepos/se-take-home-assignment

2. 点击右上角的 "Fork" 按钮

3. 选择你的账号（leonlonyio）作为 Fork 目标

4. 等待 Fork 完成


# 克隆你的 Fork（不是原仓库）
git clone https://github.com/leonlonyio/se-take-home-assignment.git

coding...

# 上传代码

git add .
git commit -m "[Backend-golang] FulianYao:McDonald's Order Controller with CLI and bot management"
git push origin main


# 添加原仓库作为远程，方便同步更新
git remote add upstream https://github.com/feedmepos/se-take-home-assignment.git

# 查看远程仓库配置
git remote -v
# 应该显示：
# origin    https://github.com/leonlonyio/se-take-home-assignment.git (fetch)
# origin    https://github.com/leonlonyio/se-take-home-assignment.git (push)
# upstream  https://github.com/feedmepos/se-take-home-assignment.git (fetch)
# upstream  https://github.com/feedmepos/se-take-home-assignment.git (push)

# 确保在 main 分支
git checkout main


# 使用 GitHub CLI
gh pr create \
  --repo feedmepos/se-take-home-assignment \
  --head leonlonyio:main \
  --base main \
  --title "[Backend-golang] FulianYao:McDonald's Order Controller with CLI and bot management" \
  --body "[Backend-golang] FulianYao:McDonald's Order Controller with CLI and bot management"" \
  --web

# 在 PR 中添加评论：
gh pr comment 1 --repo feedmepos/se-take-home-assignment --body "已完成所有功能，请审核"



