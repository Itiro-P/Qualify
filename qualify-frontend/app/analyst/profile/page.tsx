"use client";

import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import type { ApiError } from "@/libs/api";
import { Alert } from "@/components/ui";

export default function Profile() {
  const router = useRouter();
  const [analyst, setAnalyst] = useState<Analyst | null>(null);
  const [error, setError] = useState("");

  async function getAnalyst(id_user: number): Promise<Analyst> {
    return (await analystService.getByUserId(id_user)).analyst;
  }

  useEffect(() => {
    async function fetchAnalyst() {
      const session = getSessionUser();

      if (!session) {
        router.push("/user/register");
        return;
      }
      console.log(session);
      try {
        const analyst = await getAnalyst(session.id);
        setAnalyst(analyst);
      } catch (err) {
        const apiError = err as ApiError;
        if (apiError.status === 401) {
          setError("E-mail ou senha incorretos.");
        } else if (apiError.status === 400) {
          setError(apiError.message || "Dados inválidos.");
        }
      }
    }

    fetchAnalyst();
    console.log(analyst);
  }, []);

  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      {error && <Alert variant="error">{error}</Alert>}
      {analyst ? (
        <div>
          <div className="flex justify-between ml-3 mt-10">
            <ImageProfile />
            <Informations
              name={analyst.name}
              city={analyst.city}
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
        <p>Usuário não logado ou não é analista</p>
      )}

      <Footer />
    </section>
  );
}
