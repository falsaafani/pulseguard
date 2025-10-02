'use client'

import { useEffect, useState } from 'react'

interface Target {
  id: number
  name: string
  url: string
  type: string
  enabled: boolean
}

export function TargetList() {
  const [targets, setTargets] = useState<Target[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // TODO: Fetch targets from probe-api
    setLoading(false)
  }, [])

  if (loading) {
    return <div>Loading targets...</div>
  }

  return (
    <div className="space-y-4">
      {targets.length === 0 ? (
        <p className="text-gray-400">No targets configured yet.</p>
      ) : (
        targets.map((target) => (
          <div key={target.id} className="border border-gray-700 rounded p-4">
            <h3 className="font-semibold">{target.name}</h3>
            <p className="text-sm text-gray-400">{target.url}</p>
            <span className="text-xs bg-gray-800 px-2 py-1 rounded mt-2 inline-block">
              {target.type}
            </span>
          </div>
        ))
      )}
    </div>
  )
}
