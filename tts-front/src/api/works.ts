import { request } from '@/utils/request'

export interface Work {
  taskId: string
  status: string
  progress: number
  format: string
  createdAt: string
}

export interface WorksResponse {
  list: Work[]
  total: number
}

export const worksApi = {
  list() {
    return request<WorksResponse>({
      url: '/api/works/list',
      method: 'GET',
    })
  },
}
