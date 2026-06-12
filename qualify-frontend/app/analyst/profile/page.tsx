"use client";

import { Footer, Header } from "@/components";
import { useRouter, useSearchParams } from "next/navigation";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";
import { Suspense, useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import { Loading } from "@/components/ui";

function ProfileContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const profileId = searchParams.get("id")?.trim() ?? "";
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  useEffect(() => {
    async function fetchAnalyst() {
      // Visualizando o perfil de um analista específico (público, sem login).
      if (profileId) {
        const analyst = await analystService.getByUserId(Number(profileId));
        if (analyst == null) {
          router.push("/customer/searchAnalyst");
        } else {
          setAnalyst(analyst);
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
      }
    }

    fetchAnalyst();
  }, [router, profileId]);

  return analyst != null ? (
    <div>
      <div className="flex justify-between ml-3 mt-10">
        <ImageProfile />
        <Informations
          name={analyst.name}
          city={analyst.city}
          state={analyst.country_state}
          country={analyst.country_name}
          rating={analyst.mean_rating}
          reviews={analyst.total_reviews}
          date={analyst.time_created}
        />
        <ContactButtons />
      </div>

      <div className="flex mt-10">
        <About analyst={analyst} />
        <StatisticProfile
          analyst_id={analyst.id}
          hourly_rate={analyst.hourly_rate}
        />
      </div>
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
