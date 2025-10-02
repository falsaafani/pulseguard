'use client'

import { useEffect, useState } from 'react'

interface Incident {
  id: number
  target_id: number
  started_at: string
  ended_at?: string
  kind: string
  details: string
}

export function IncidentTimeline() {
  const [incidents, setIncidents] = useState<Incident[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // TODO: Fetch incidents from probe-api or database
    setLoading(false)
  }, [])

  if (loading) {
    return <div>Loading incidents...</div>
  }

  return (
    <div className="space-y-4">
      {incidents.length === 0 ? (
        <p className="text-gray-400">No incidents recorded.</p>
      ) : (
        incidents.map((incident) => (
          <div key={incident.id} className="border-l-4 border-red-500 pl-4 py-2">
            <p className="font-semibold">{incident.kind}</p>
            <p className="text-sm text-gray-400">{incident.details}</p>
            <p className="text-xs text-gray-500 mt-1">
              {new Date(incident.started_at).toLocaleString()}
            </p>
          </div>
        ))
      )}
    </div>
  )
}
