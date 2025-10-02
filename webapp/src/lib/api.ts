const API_BASE_URL = process.env.NEXT_PUBLIC_PROBE_API_URL || 'http://localhost:8080'

export async function fetchTargets() {
  const response = await fetch(`${API_BASE_URL}/targets`)
  if (!response.ok) {
    throw new Error('Failed to fetch targets')
  }
  return response.json()
}

export async function fetchStatus(targetId?: number) {
  const url = targetId
    ? `${API_BASE_URL}/status?target_id=${targetId}`
    : `${API_BASE_URL}/status`

  const response = await fetch(url)
  if (!response.ok) {
    throw new Error('Failed to fetch status')
  }
  return response.json()
}

export async function createTarget(name: string, url: string, type: string) {
  const response = await fetch(`${API_BASE_URL}/targets`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, url, type }),
  })

  if (!response.ok) {
    throw new Error('Failed to create target')
  }
  return response.json()
}
