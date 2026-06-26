"use client";

import { Footer, Header } from "@/components";
import { useRouter, useSearchParams } from "next/navigation";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
  ProposalModal,
} from "@/components/analyst/profile";
import { Suspense, useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService, clientService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import { Loading } from "@/components/ui";

function ProfileContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const profileId = searchParams.get("id")?.trim() ?? "";
  const [analyst, setAnalyst] = useState<Analyst | null>(null);
  const [picture, setPicture] = useState<string>("");
  const [clientId, setClientId] = useState<number | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [proposalSent, setProposalSent] = useState(false);

  useEffect(() => {
    async function fetchAnalyst() {
      // Visualizando o perfil de um analista específico (público, sem login).
      if (profileId) {
        const analyst = await analystService.getByUserId(Number(profileId));
        if (analyst == null) {
          router.push("/customer/searchAnalyst");
        } else {
          setAnalyst(analyst);
          const profileResp = await analystService.getProfile(analyst.id);
          setPicture(profileResp?.picture ?? "");
        }

        // Tenta carregar o client do usuário logado para habilitar o botão.
        const session = await getSessionUser();
        if (session) {
          const client = await clientService.getByUserId(session.id);
          if (client?.id) setClientId(client.id);
        }
        return;
      }

      // Sem id: é o "meu perfil" — exige autenticação.
      const session = await getSessionUser();
      if (!session) {
        router.push("/user/login");
        return;
      }

      const analyst = await analystService.getByUserId(session.id);
      if (analyst == null) {
        router.push("/analyst/register");
      } else {
        setAnalyst(analyst);
        const profileResp = await analystService.getProfile(analyst.id);
        setPicture(profileResp?.picture ?? "");
      }
    }

    fetchAnalyst();
  }, [router, profileId]);

  const isPublicView = profileId !== "";
  const canPropose = isPublicView && clientId !== null && !proposalSent;

  return analyst != null ? (
    <div>
      <div className="flex justify-between ml-3 mt-10">
        <ImageProfile user={analyst.name} imageURL={picture} />
        <Informations
          name={analyst.name}
          city={analyst.city}
          state={analyst.country_state}
          country={analyst.country_name}
          rating={analyst.mean_rating}
          reviews={analyst.total_reviews}
          date={analyst.time_created}
        />
        <ContactButtons
          showPropose={canPropose}
          onPropose={() => setShowModal(true)}
        />
      </div>

      {proposalSent && (
        <div className="mt-4 ml-3 rounded-lg border border-accent/30 bg-accent/10 px-4 py-3 text-sm text-accent">
          Proposta enviada com sucesso!
        </div>
      )}

      <div className="flex mt-10">
        <About analyst={analyst} />
        <StatisticProfile
          analyst_id={analyst.id}
          hourly_rate={analyst.hourly_rate}
        />
      </div>

      {showModal && analyst && clientId && (
        <ProposalModal
          analystId={analyst.id}
          clientId={clientId}
          onClose={() => setShowModal(false)}
          onSuccess={() => {
            setShowModal(false);
            setProposalSent(true);
          }}
        />
      )}
    </div>
  ) : (
    <Loading />
  );
}

export default function Profile() {
  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      <Suspense fallback={<Loading />}>
        <ProfileContent />
      </Suspense>
      <Footer />
    </section>
  );
}
