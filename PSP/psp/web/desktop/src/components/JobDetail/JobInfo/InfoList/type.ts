export enum InfoBlockType {
  Normal,
  Other,
}
export enum FieldType {
  Normal,
  File,
}
export interface InfoBlockTemplate {
  title: string
  icon: string
  children: Field[]
  partitionNum?: number
  type?: InfoBlockType
}

export interface Field {
  key: string
  text: string
  share?: number
  type?: FieldType
}
