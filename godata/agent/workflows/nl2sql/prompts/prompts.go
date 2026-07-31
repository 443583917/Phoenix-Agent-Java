package prompts

import _ "embed"

// ──────────────────────────── Intent Recognition ────────────────────────────

// IntentRecognitionPrompt classifies the user's query as either casual chat or
// a potential data analysis request.
const IntentRecognitionPrompt = `你是一个意图识别助手。请判断用户的问题是"闲聊或无关指令"还是"可能的数据分析请求"。

判断标准：
- "闲聊或无关指令"：纯粹的社交对话、情感表达、无关话题、简单问候、无关指令等
- "可能的数据分析请求"：涉及数据查询、报表、统计、分析、指标、趋势、对比、数据库等数据相关请求

请只输出一个JSON对象，格式如下：
{"classification": "闲聊或无关指令", "confidence": 0.95, "reasoning": "原因"}

用户问题：{{.Query}}

历史上下文：{{.Context}}`

// ──────────────────────────── Evidence Rewrite ────────────────────────────

// EvidenceRewritePrompt rewrites the query for evidence retrieval.
const EvidenceRewritePrompt = `你是一个查询改写助手。请将用户的问题改写为适合进行知识检索的独立查询。

要求：
1. 解析指代：将"它"、"它们"、"这个"、"那个"等代词替换为具体内容
2. 补全上下文：将省略的部分补全，使其成为一个完整的、独立的问题
3. 保持原意：不要改变用户的核心意图
4. 输出JSON格式：{"canonical_query": "改写后的查询"}

当前对话上下文：
{{.Context}}

用户当前问题：{{.Query}}`

// ──────────────────────────── Query Enhancement ────────────────────────────

// QueryEnhancementPrompt normalizes and expands the user's query.
const QueryEnhancementPrompt = `你是一个查询增强助手。请对用户的数据分析问题进行规范化，并生成多个扩展查询以便全面检索。

要求：
1. 规范化查询：修正拼写、补全缩写、统一术语
2. 生成扩展查询：从不同角度生成2-4个相关查询，便于知识库检索
3. 输出JSON格式：
{"canonical_query": "规范化后的主要查询", "expanded_queries": ["扩展查询1", "扩展查询2", "扩展查询3"]}

用户问题：{{.Query}}

改写后的检索查询：{{.RewrittenQuery}}`

// ──────────────────────────── Feasibility Assessment ────────────────────────────

// FeasibilityAssessmentPrompt classifies the query type for routing.
const FeasibilityAssessmentPrompt = `你是一个需求可行性评估助手。请分析用户的数据分析请求，判断其类型。

分类选项：
1. "数据分析" — 明确的数据分析需求，可以通过SQL查询和可视化完成（如：查询销售额、统计用户数、分析趋势等）
2. "需要澄清" — 需求模糊、缺少关键信息，需要进一步询问（如：笼统地问"分析一下数据"）
3. "自由闲聊" — 非数据分析相关的自由对话（如：问候、闲聊、问问题等）

输出JSON格式：
{"classification": "数据分析", "reasoning": "简洁的原因说明"}

用户问题：{{.Query}}

可用的数据表信息：
{{.SchemaInfo}}`

// ──────────────────────────── Planner ────────────────────────────

// PlannerPrompt generates a multi-step execution plan.
const PlannerPrompt = `你是一个数据分析规划助手。请为以下数据分析请求生成一个执行计划。

可用的工具：
1. sql_generate — 生成SQL查询语句
2. python_execute — 执行Python代码进行数据分析

请输出一个JSON格式的执行计划：
{
  "thought_process": "思考过程的简要描述",
  "execution_plan": [
    {"step": 1, "tool_to_use": "sql_generate", "tool_parameters": {"instruction": "步骤说明"}},
    {"step": 2, "tool_to_use": "python_execute", "tool_parameters": {"instruction": "步骤说明"}}
  ]
}

用户问题：{{.Query}}

数据表关系信息：
{{.TableRelation}}

对话上下文：
{{.Context}}`

// ──────────────────────────── SQL Generate ────────────────────────────

// SQLGeneratePrompt generates SQL from schema and evidence.
const SQLGeneratePrompt = `你是一个SQL生成助手。请根据数据库表结构和用户需求生成SQL查询语句。

数据库方言：{{.Dialect}}

表结构信息：
{{.SchemaInfo}}

相关业务知识：
{{.Evidence}}

用户需求：{{.Query}}

要求：
1. 只生成SELECT查询语句
2. 使用合适的JOIN、WHERE、GROUP BY等
3. 确保SQL语法正确
4. 输出JSON格式：{"sql": "生成的SQL语句", "explanation": "SQL说明"}

{{if .RetryReason}}之前的SQL被拒绝，原因：{{.RetryReason}}
请根据原因修正SQL。{{end}}`

// ──────────────────────────── Semantic Consistency ────────────────────────────

// SemanticConsistencyPrompt validates SQL against the user's query.
const SemanticConsistencyPrompt = `你是一个SQL语义校验助手。请检查生成的SQL是否准确地满足了用户的需求。

用户需求：{{.Query}}

生成的SQL：
{{.SQL}}

业务知识：
{{.Evidence}}

表结构：
{{.SchemaInfo}}

请验证：
1. SQL是否正确地反映了用户的查询意图
2. 是否使用了正确的表和字段
3. 聚合条件、过滤条件是否准确
4. 是否遗漏了任何用户要求的信息

输出JSON格式：
{"result": "通过", "reason": "验证说明"}

如果SQL有问题，输出：
{"result": "不通过", "reason": "具体问题说明"}`

// ──────────────────────────── Python Generate ────────────────────────────

// PythonGeneratePrompt generates Python code for data analysis.
const PythonGeneratePrompt = `你是一个Python数据分析代码生成助手。请根据SQL查询结果生成Python数据分析代码。

可用的库：pandas, numpy, matplotlib, seaborn

SQL查询结果说明：
{{.SQLResult}}

用户需求：{{.Query}}

生成说明：{{.Instruction}}

要求：
1. 代码必须完整可运行
2. 使用pandas读取数据
3. 生成必要的图表（中文标签）
4. 输出分析结论
5. 输出JSON格式：
{"code": "Python代码", "explanation": "代码说明"}

{{if .RetryInfo}}前一次代码问题：{{.RetryInfo}}
请修正。{{end}}`

// ──────────────────────────── Report Generate ────────────────────────────

// ReportGeneratePrompt generates a Markdown report.
const ReportGeneratePrompt = `你是一个数据报告生成助手。请根据分析结果生成Markdown格式的数据分析报告。

用户需求：{{.Query}}

SQL查询：{{.SQL}}

查询结果：{{.SQLResult}}

Python分析结果：{{.PythonAnalysis}}

要求：
1. 使用Markdown格式
2. 包含数据概览、分析结论、数据可视化说明
3. 语言通顺、结构清晰
4. 如果有图表，使用![]()引用

{{if .Continuation}}这是续写部分，请根据前面的内容继续完善报告。{{end}}`

// ──────────────────────────── Table Relation ────────────────────────────

// TableRelationPrompt analyzes table relationships and generates semantic model descriptions.
const TableRelationPrompt = `你是表关系分析专家。请分析以下数据库表结构信息，识别表之间的逻辑关系（外键、关联关系）。

表结构信息：
{{.SchemaInfo}}

数据库方言：{{.DBDialect}}

请返回JSON格式：
{
    "tableRelationDescription": "表之间的逻辑关系描述",
    "semanticModelPrompt": "基于表结构生成的语义模型描述"
}`

// ──────────────────────────── Python Analyze ────────────────────────────

// PythonAnalyzePrompt generates data analysis reports from SQL and Python results.
const PythonAnalyzePrompt = `你是数据分析专家。请根据以下信息生成分析报告：

SQL查询结果：
%s

Python执行结果：
%s

请分析数据并给出关键洞察和建议。返回JSON格式：
{
    "analysis": "分析结果",
    "keyInsights": ["关键洞察1", "关键洞察2"],
    "recommendations": ["建议1", "建议2"]
}`
