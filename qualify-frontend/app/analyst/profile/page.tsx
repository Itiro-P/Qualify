"use client";

import { Footer, Header } from "@/components";
import {
  ImageProfile,
  StatisticProfile,
  About,
  Informations,
  ContactButtons,
} from "@/components/analyst/profile";

import { getSessionUser } from "@/libs/session";
import { useState, useEffect } from "react";
import { Analyst } from "@/types/services";
import { analystService } from "@/libs";

export default function Profile() {
  const [analyst, setAnalyst] = useState<Analyst | null>(null);

  useEffect(() => {
    const session = getSessionUser();

    if (!session) {
      window.location.href = "/User/register";
      return;
    }

    analystService
      .getByUserId(session.id)
      .then((resp) => {
        console.log(resp);
        setAnalyst(resp.analyst);
      })
      .catch((error) => {
        console.error("ERRO:", error);

        if (error.response) {
          console.error(error.response.data);
        }
      });
  }, []);

  return (
    <section id="profile" className="px-6 md:px-20 py-14">
      <Header />
      {analyst ? (
        <div>
          <div className="flex justify-between ml-3 mt-10">
            <ImageProfile />
            <Informations />
            <ContactButtons />
          </div>

          <div className="flex mt-10">
            <About analyst={analyst} />
            <StatisticProfile />
          </div>
        </div>
      ) : (
        <p>Usuário não logado ou não é analista</p>
      )}

      <Footer />
    </section>
  );
}
