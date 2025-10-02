import { TargetList } from '@/components/TargetList'
import { IncidentTimeline } from '@/components/IncidentTimeline'

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <h1 className="text-4xl font-bold mb-8">PulseGuard Dashboard</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <section>
          <h2 className="text-2xl font-semibold mb-4">Monitored Targets</h2>
          <TargetList />
        </section>

        <section>
          <h2 className="text-2xl font-semibold mb-4">Recent Incidents</h2>
          <IncidentTimeline />
        </section>
      </div>
    </main>
  )
}
