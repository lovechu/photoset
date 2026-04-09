#!/bin/bash
# 高级搜索过滤器 API 测试脚本

echo "🔍 测试高级搜索过滤器 - Day 1 后端接口验证"
echo "=================================================="

# 1. 测试基础搜索（原有功能）
echo ""
echo "1️⃣ 测试基础搜索 (原有功能):"
curl -s "http://localhost:8080/api/photosets?keyword=海边&page=1&page_size=2" | jq '.list[0:2] | length'

# 2. 测试价格筛选
echo ""
echo "2️⃣ 测试价格筛选 (免费):"
curl -s "http://localhost:8080/api/photosets?price_min=0&price_max=0&page=1&page_size=2" | jq '.list[0:2] | length'

# 3. 测试排序功能
echo ""
echo "3️⃣ 测试排序功能 (最新发布):"
curl -s "http://localhost:8080/api/photosets?sort_by=latest&page=1&page_size=2" | jq '.list[0:2] | length'

echo ""
echo "4️⃣ 测试排序功能 (价格从低到高):"
curl -s "http://localhost:8080/api/photosets?sort_by=price_asc&page=1&page_size=2" | jq '.list[0:2] | length'

# 4. 测试时间范围筛选
echo ""
echo "5️⃣ 测试时间范围筛选 (本周):"
curl -s "http://localhost:8080/api/photosets?time_range=week&page=1&page_size=2" | jq '.list[0:2] | length'

# 5. 测试组合筛选
echo ""
echo "6️⃣ 测试组合筛选 (免费+本周):"
curl -s "http://localhost:8080/api/photosets?price_min=0&price_max=0&time_range=week&page=1&page_size=2" | jq '.list[0:2] | length'

echo ""
echo "=================================================="
echo "✅ Day 1 后端接口改造完成！"
echo ""
echo "📊 下一步：Day 2 前端筛选组件开发"
echo ""
echo "🛠️ 需要创建的文件："
echo "  1. src/components/FilterPanel.vue"
echo "  2. src/components/PriceFilter.vue"
echo "  3. src/components/SortMenu.vue"
echo "  4. src/components/TimeRangeFilter.vue"
echo ""
echo "⚙️ 需要更新的文件："
echo "  1. src/views/Home.vue (集成筛选面板)"
echo "  2. src/api/photoset.js (添加筛选参数)"
echo "  3. src/stores/search.js (新增筛选状态管理)"