// 检查sqlite是否存在
if [ ! -f /data/tokens.db ]; then
    # 如果不存在则创建
    gitea-matrix-bot initdb 
fi
# 启动服务
gitea-matrix-bot -v 
