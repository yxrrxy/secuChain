import json
import csv

def parse_cyclonedx_bom(file_path, output_csv):
    """
    从指定的 CycloneDX BOM 文件中解析构件信息，并生成 CSV 文件。
    :param file_path: CycloneDX BOM 文件路径
    :param output_csv: 输出的 CSV 文件路径
    """
    try:
        # 读取 CycloneDX BOM 文件
        with open(file_path, 'r', encoding='utf-8') as file:
            bom_data = json.load(file)

        # 提取构件信息
        components = bom_data.get('components', [])

        # 定义 CSV 文件的表头
        csv_headers = ['构件名称', '版本', '类型', 'PURL']

        # 写入 CSV 文件，使用 utf-8-sig 编码
        with open(output_csv, 'w', newline='', encoding='utf-8-sig') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(csv_headers)  # 写入表头
            for component in components:
                name = component.get('name', '')
                version = component.get('version', '')
                component_type = component.get('type', '')
                purl = component.get('purl', '')
                writer.writerow([name, version, component_type, purl])

        print(f"成功生成 CSV 文件：{output_csv}")
    except FileNotFoundError:
        print(f"错误：文件 {file_path} 未找到。")
    except json.JSONDecodeError:
        print(f"错误：文件 {file_path} 不是有效的 JSON 格式。")
    except Exception as e:
        print(f"发生错误：{e}")

# 示例用法
bom_file_path = 'output.cdx.json'  # 替换为你的 CycloneDX BOM 文件路径
output_csv_file = 'components_table.csv'  # 输出的 CSV 文件路径

parse_cyclonedx_bom(bom_file_path, output_csv_file)