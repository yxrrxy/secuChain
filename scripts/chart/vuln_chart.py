import json
import matplotlib.pyplot as plt


#将json转换为一定格式
# 读取JSON文件
file_path = 'CVE-2025.json'  
with open(file_path, 'r', encoding='utf-8') as file:
    data = json.load(file)

# 提取漏洞来源和安全等级数据
source_data = {}
severity_data = {}

for item in data['cve_items']:
    # 提取来源
    source = item['sourceIdentifier'].split('@')[-1]  # 提取@之后的部分，如果没有@则取整个字符串

    # 提取安全等级
    if 'cvssMetricV31' in item['metrics'] and item['metrics']['cvssMetricV31']:
        severity = item['metrics']['cvssMetricV31'][0]['cvssData']['baseSeverity']
    else:
        severity = 'Unknown'  # 如果没有CVSS数据，标记为Unknown

    # 统计来源
    if source in source_data:
        source_data[source] += 1
    else:
        source_data[source] = 1

    # 统计安全等级
    if severity in severity_data:
        severity_data[severity] += 1
    else:
        severity_data[severity] = 1

# 创建柱状图（漏洞来源）
sources = list(source_data.keys())
source_counts = list(source_data.values())

plt.figure(figsize=(10, 6))
plt.bar(sources, source_counts, color=['blue', 'green', 'red', 'purple', 'orange'])
plt.title("Vulnerability Sources")
plt.xlabel("Source")
plt.ylabel("Number of Vulnerabilities")
for i, v in enumerate(source_counts):
    plt.text(i, v + 1, str(v), ha='center')  # 在柱状图上显示数值
plt.xticks(rotation=45)  # 如果来源名称过长，可以旋转标签
plt.tight_layout()  # 自动调整布局

# 保存柱状图为文件
plt.savefig('vulnerability_sources.png')  # 保存为PNG文件
plt.close()  # 关闭当前图表

# 创建饼状图（漏洞安全性）
labels = list(severity_data.keys())
values = list(severity_data.values())

plt.figure(figsize=(8, 8))
plt.pie(values, labels=labels, autopct='%1.1f%%', startangle=140, colors=['red', 'orange', 'yellow', 'green', 'gray'])
plt.title("Vulnerability Severity Distribution")

# 保存饼状图为文件
plt.savefig('vulnerability_severity.png')  # 保存为PNG文件
plt.close()  # 关闭当前图表

print("Charts have been saved as 'vulnerability_sources.png' and 'vulnerability_severity.png'.")