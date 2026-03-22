import { SectionCards } from "../components/section-cards";

export default function Dashboard() {
  return (
    <div className="flex flex-col gap-6 bg-black text-white">
      <div className="px-4 lg:px-6">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground text-sm">
          Crixus PLC Payroll & Banking System
        </p>
      </div>
      <SectionCards />
    </div>
  );
}
