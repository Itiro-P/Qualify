"use client";

import { Loading } from "@/components/ui/Loading";
import { ProposalCard } from "@/components/analyst/proposals";
import { analystService, proposalService } from "@/libs/services";
import { serviceService } from "@/libs/services/serviceService";
import { User } from "@/libs/session";
import { ProposalLetter } from "@/types/services/proposal";
import { Service } from "@/types/services/service";
import { Briefcase, Clock } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";

type StatusKey = "PROPOSTA" | "ATIVO" | "FINALIZADO" | "CANCELADO";
type ServiceRole = "client" | "analyst";

interface ServiceListItem {
  service: Service;
  role: ServiceRole;
  statusKey: Exclude<StatusKey, "PROPOSTA">;
}

const TABS: { key: StatusKey; label: string }[] = [
  { key: "PROPOSTA", label: "Propostas" },
  { key: "ATIVO", label: "Ativos" },
  { key: "FINALIZADO", label: "Finalizados" },
  { key: "CANCELADO", label: "Cancelados" },
];

const STATUS_BADGE: Record<StatusKey, string> = {
  PROPOSTA: "bg-blue-500/10 text-blue-400 border-blue-500/30",
  ATIVO: "bg-accent/10 text-accent border-accent/30",
  FINALIZADO: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
  CANCELADO: "bg-red-500/10 text-red-400 border-red-500/30",
};

function normalizeStatus(value?: string): Exclude<StatusKey, "PROPOSTA"> {
  const base = (value ?? "")
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "");

  if (base.includes("cancel") || base.includes("block")) {
    return "CANCELADO";
  }

  if (
    base.includes("final") ||
    base.includes("conclu") ||
    base.includes("complet")
  ) {
    return "FINALIZADO";
  }

  return "ATIVO";
}

function ServiceCard({
  item,
  acting,
  onCancel,
  onFinalize,
}: {
  item: ServiceListItem;
  acting: boolean;
  onCancel: () => void;
  onFinalize: () => void;
}) {
  const { service, statusKey, role } = item;
  const label = TABS.find((t) => t.key === statusKey)?.label ?? statusKey;
  const canCancel =
    role === "analyst" && statusKey !== "CANCELADO" && statusKey !== "FINALIZADO";
  const canFinalize =
    role === "client" && statusKey !== "CANCELADO" && statusKey !== "FINALIZADO";
  const contextLabel = role === "analyst" ? "Como analista" : "Como cliente";

  return (
    <div className="rounded-xl border border-white/10 bg-white/5 p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h3 className="text-lg font-bold text-white">{service.title}</h3>
          <p className="mt-1 text-xs text-neutral-slate">{contextLabel}</p>
        </div>
        <span className={`shrink-0 rounded-full border px-3 py-1 text-xs font-medium ${STATUS_BADGE[statusKey]}`}>
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

      {canCancel && (
        <div className="mt-5">
          <button
            onClick={onCancel}
            disabled={acting}
            className="w-full rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-2 text-sm font-medium hover:bg-red-500/20 transition-colors disabled:opacity-50 cursor-pointer"
          >
            Cancelar serviço
          </button>
        </div>
      )}

      {canFinalize && (
        <div className="mt-5">
          <button
            onClick={onFinalize}
            disabled={acting}
            className="w-full rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 px-4 py-2 text-sm font-medium hover:bg-emerald-500/20 transition-colors disabled:opacity-50 cursor-pointer"
          >
            Finalizar serviço
          </button>
        </div>
      )}
    </div>
  );
}

export function ListServices({ user }: { user: User }) {
  const [services, setServices] = useState<ServiceListItem[]>([]);
  const [proposals, setProposals] = useState<ProposalLetter[]>([]);
  const [loading, setLoading] = useState(false);
  const [actingServiceId, setActingServiceId] = useState<number | null>(null);
  const [activeTab, setActiveTab] = useState<StatusKey>("ATIVO");

  const loadData = useCallback(async () => {
    setLoading(true);

    const [clientServices, analyst] = await Promise.all([
      serviceService.listServicesByClient(user.id),
      analystService.getByUserId(user.id),
    ]);

    let analystServices: Service[] = [];
    let pendingProposals: ProposalLetter[] = [];

    if (analyst) {
      const [servicesByAnalyst, analystProposals] = await Promise.all([
        analystService.listServices(user.id),
        proposalService.listByAnalyst(user.id),
      ]);

      analystServices = servicesByAnalyst ?? [];

      const pendingChecks = await Promise.all(
        (analystProposals ?? []).map(async (proposal) => {
          if (!proposal.id) return proposal;
          const linkedServices = await serviceService.list({
            proposal_id: proposal.id,
          });

          return linkedServices && linkedServices.length > 0 ? null : proposal;
        }),
      );

      pendingProposals = pendingChecks.filter(
        (proposal): proposal is ProposalLetter => proposal != null,
      );
    }

    const clientList: ServiceListItem[] = (clientServices ?? []).map((service) => ({
      service,
      role: "client",
      statusKey: normalizeStatus(service.status),
    }));

    const analystList: ServiceListItem[] = analystServices.map((service) => ({
      service,
      role: "analyst",
      statusKey: normalizeStatus(service.status),
    }));

    setServices([...clientList, ...analystList]);
    setProposals(pendingProposals);
    setLoading(false);
  }, [user.id]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  async function handleCancelService(serviceId?: number) {
    if (!serviceId) return;
    setActingServiceId(serviceId);
    await serviceService.patch(serviceId, { status: "Cancelado" });
    setActingServiceId(null);
    await loadData();
  }

  async function handleFinalizeService(serviceId?: number) {
    if (!serviceId) return;
    setActingServiceId(serviceId);
    await serviceService.patch(serviceId, { status: "Finalizado" });
    setActingServiceId(null);
    await loadData();
  }

  const countByStatus = useMemo(() => {
    const counts = {} as Record<StatusKey, number>;
    counts.PROPOSTA = proposals.length;
    counts.ATIVO = services.filter((s) => s.statusKey === "ATIVO").length;
    counts.FINALIZADO = services.filter((s) => s.statusKey === "FINALIZADO").length;
    counts.CANCELADO = services.filter((s) => s.statusKey === "CANCELADO").length;
    return counts;
  }, [services, proposals]);

  const visibleServices = useMemo(() => {
    if (activeTab === "PROPOSTA") return [];
    return services.filter((s) => s.statusKey === activeTab);
  }, [services, activeTab]);

  if (loading) return <Loading />;

  return (
    <section className="px-6 md:px-20 py-14">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">Meus serviços</h1>
          <div className="h-1 w-20 bg-accent mb-6" />
          <p className="text-neutral-slate max-w-2xl">
            Acompanhe propostas e serviços por status.
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

        {activeTab === "PROPOSTA" ? (
          proposals.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
              <Briefcase className="w-10 h-10 text-neutral-slate mb-4" />
              <p className="text-neutral-slate">Nenhuma proposta pendente.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {proposals.map((proposal) => (
                <ProposalCard
                  key={proposal.id}
                  proposal={proposal}
                  viewerRole="analyst"
                  onUpdate={() => {
                    void loadData();
                  }}
                />
              ))}
            </div>
          )
        ) : visibleServices.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 bg-white/2 py-20 text-center">
            <Briefcase className="w-10 h-10 text-neutral-slate mb-4" />
            <p className="text-neutral-slate">
              Nenhum serviço nesta etapa.
            </p>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {visibleServices.map((item) => (
              <ServiceCard
                key={`${item.role}-${item.service.id ?? item.service.title}`}
                item={item}
                acting={actingServiceId === item.service.id}
                onCancel={() => {
                  void handleCancelService(item.service.id);
                }}
                onFinalize={() => {
                  void handleFinalizeService(item.service.id);
                }}
              />
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
