export type DataType = {
  group: string
  cols: DataColumn[]
}

export type DataColumn = {
  X: string
  Y: {
    members: number
    percent: number
  }
}
