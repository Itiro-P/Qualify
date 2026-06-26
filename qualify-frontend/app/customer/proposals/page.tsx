"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Header, Footer } from "@/components";
import { Loading } from "@/components/ui";
import { ProposalCard } from "@/components/analyst/proposals";
import { proposalService } from "@/libs/services";
import { getSessionUser } from "@/libs/session";
import type { ProposalLetter } from "@/types/services/proposal";

export default function CustomerProposalsPage() {
  const router = useRouter();
  const [proposals, setProposals] = useState<ProposalLetter[]>([]);
  const [loading, setLoading] = useState(true);
  const [userId, setUserId] = useState<number | null>(null);

  const loadProposals = useCallback(
    async (uid: number) => {
      const list = await proposalService.listByClient(uid);
      setProposals(list ?? []);
      setLoading(false);
    },
    [],
  );

  useEffect(() => {
    getSessionUser().then((user) => {
      if (!user) {
        router.push("/user/login");
        return;
      }
      setUserId(user.id);
      loadProposals(user.id);
    });
  }, [router, loadProposals]);

  if (loading) {
    return <Loading />;
  }

  return (
    <section className="px-6 md:px-20 py-14 min-h-screen">
      <Header />

      <div className="mt-10">
        <h1 className="text-2xl font-bold text-white mb-2">
          Minhas Propostas
        </h1>
        <p className="text-sm text-neutral-slate mb-8">
          Acompanhe as propostas que você enviou para analistas de QA.
        </p>

        {proposals.length === 0 ? (
          <div className="rounded-xl border border-white/10 bg-white/5 p-10 text-center">
            <p className="text-neutral-slate">
              Nenhuma proposta enviada ainda.
            </p>
            <button
              onClick={() => router.push("/customer/searchAnalyst")}
              className="mt-4 text-sm text-accent hover:underline"
            >
              Encontrar um analista
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {proposals.map((proposal) => (
              <ProposalCard
                key={proposal.id}
                proposal={proposal}
                viewerRole="client"
                onUpdate={() => userId && loadProposals(userId)}
              />
            ))}
          </div>
        )}
      </div>

      <Footer />
    </section>
  );
}
