"use client";

import { Loading } from "@/components/ui/Loading";
import { ProposalCard } from "@/components/analyst/proposals";
import { analystService, proposalService } from "@/libs/services";
import { serviceService } from "@/libs/services/serviceService";
import { User } from "@/libs/session";
import { ProposalLetter } from "@/types/services/proposal";
import { Service } from "@/types/services/service";
import { Briefcase, Clock } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

type StatusKey =
  | "PROPOSAL"
  | "NEGOTIATION"
  | "IN_PROGRESS"
  | "BLOCKED"
  | "COMPLETED";

const TABS: { key: StatusKey; label: string }[] = [
  { key: "PROPOSAL", label: "Propostas" },
  { key: "NEGOTIATION", label: "Negociação" },
  { key: "IN_PROGRESS", label: "Em andamento" },
  { key: "BLOCKED", label: "Bloqueados" },
  { key: "COMPLETED", label: "Concluídos" },
];

const STATUS_BADGE: Record<StatusKey, string> = {
  PROPOSAL: "bg-blue-500/10 text-blue-400 border-blue-500/30",
  NEGOTIATION: "bg-yellow-500/10 text-yellow-400 border-yellow-500/30",
  IN_PROGRESS: "bg-accent/10 text-accent border-accent/30",
  BLOCKED: "bg-red-500/10 text-red-400 border-red-500/30",
  COMPLETED: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
};

function ServiceCard({ service }: { service: Service }) {
  const status = (service.status as StatusKey) ?? "PROPOSAL";
  const label = TABS.find((t) => t.key === status)?.label ?? status;

  return (
    <div className="rounded-xl border border-white/10 bg-white/5 p-6">
      <div className="flex items-start justify-between gap-4">
        <h3 className="text-lg font-bold text-white">{service.title}</h3>
        <span
          className={`shrink-0 rounded-full border px-3 py-1 text-xs font-medium ${
            STATUS_BADGE[status] ?? STATUS_BADGE.PROPOSAL
          }`}
        >
          {label}
        </span>
      </div>
      <p className="mt-3 text-sm leading-relaxed text-slate-300">
        {service.content}
      </p>
      {service.hourly_rate != null && (
        <p className="mt-4 flex items-center gap-1.5 text-sm text-neutral-slate">
          <Clock className="w-4 h-4 text-accent" />
          R$ {service.hourly_rate}/h
        </p>
      )}
    </div>
  );
}

export function ListServices({ user }: { user: User }) {
  const [services, setServices] = useState<Service[]>([]);
  const [proposals, setProposals] = useState<ProposalLetter[]>([]);
  const [isAnalyst, setIsAnalyst] = useState(false);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<StatusKey>("COMPLETED");

  useEffect(() => {
    async function getInfo() {
      setLoading(true);
      const analyst = await analystService.getByUserId(user.id);

      if (analyst) {
        setIsAnalyst(true);
        const analystProposals = await proposalService.listByAnalyst(user.id);
        setProposals(analystProposals ?? []);
        setServices([]);
      } else {
        setIsAnalyst(false);
        const services = await serviceService.listServicesByClient(user.id);
        if (services) {
          setServices(services);
        } else {
          setServices([]);
        }
        setProposals([]);
      }

      setLoading(false);
    }

    getInfo();
  }, [user.id]);

  const countByStatus = useMemo(() => {
    const counts = {} as Record<StatusKey, number>;
    for (const tab of TABS) {
      counts[tab.key] = services.filter((s) => s.status === tab.key).length;
    }
    return counts;
  }, [services]);

  const visibleServices = useMemo(
    () => services.filter((s) => s.status === activeTab),
    [services, activeTab],
  );

  if (loading) return <Loading />;

  if (isAnalyst) {
    return (
      <section className="px-6 md:px-20 py-14">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-3xl font-bold mb-4">Meus serviços</h1>
            <div className="h-1 w-20 bg-accent mb-6" />
            <p className="text-neutral-slate max-w-2xl">
              Aceite ou recuse propostas recebidas de clientes.
            </p>
          </div>

          {proposals.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
              <Briefcase className="w-10 h-10 text-neutral-slate mb-4" />
              <p className="text-neutral-slate">Nenhuma proposta recebida.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {proposals.map((proposal) => (
                <ProposalCard
                  key={proposal.id}
                  proposal={proposal}
                  viewerRole="analyst"
                  onUpdate={async () => {
                    const updatedProposals = await proposalService.listByAnalyst(
                      user.id,
                    );
                    setProposals(updatedProposals ?? []);
                  }}
                />
              ))}
            </div>
          )}
        </div>
      </section>
    );
  }

  return (
    <section className="px-6 md:px-20 py-14">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">Serviços</h1>
          <div className="h-1 w-20 bg-accent mb-6" />
          <p className="text-neutral-slate max-w-2xl">
            Acompanhe seus serviços organizados por etapa.
          </p>
        </div>

        <div className="flex flex-wrap gap-2 border-b border-white/10 mb-6">
          {TABS.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 -mb-px transition-colors cursor-pointer ${
                activeTab === tab.key
                  ? "border-accent text-accent"
                  : "border-transparent text-neutral-slate hover:text-white"
              }`}
            >
              {tab.label}
              <span className="rounded-full bg-white/10 px-2 py-0.5 text-xs">
                {countByStatus[tab.key]}
              </span>
            </button>
          ))}
        </div>

        {visibleServices.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
            <Briefcase className="w-10 h-10 text-neutral-slate mb-4" />
            <p className="text-neutral-slate">
              Nenhum serviço nesta etapa.
            </p>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {visibleServices.map((service) => (
              <ServiceCard key={service.id} service={service} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
