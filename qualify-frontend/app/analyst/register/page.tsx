"use client";

import { RegisterAnalyst } from "@/components/analyst/register";
import { Footer, Header } from "@/components";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { userService } from "@/libs";
import { getSessionUser } from "@/libs/session";
import { Loading } from "@/components/ui";
import { User } from "@/types/services";

export default function Register() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function fetchAnalyst() {
      const session = getSessionUser();

      if (!session) {
        router.push("/user/register");
        return;
      }
      const user = await userService.getById(session.id);
      if (user == null) {
        router.push("/user/login");
      } else {
        setUser(user);
      }
    }

    fetchAnalyst();
  }, [router]);
  return (
    <section id="register" className="px-6 md:px-20 py-14">
      <Header />
      <div className="flex justify-center">
        <div className="flex flex-col ml-3 w-1/2">
          {user != null ? <RegisterAnalyst userSession={user} /> : <Loading />}
        </div>
      </div>
      <Footer />
    </section>
  );
}
