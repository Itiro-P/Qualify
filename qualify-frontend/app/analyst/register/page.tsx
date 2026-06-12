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
    async function fetchUser() {
      const session = await getSessionUser();

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

    fetchUser();
  }, [router]);
  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      <Header />
      <main className="flex-1">
        {user != null ? <RegisterAnalyst userSession={user} /> : <Loading />}
      </main>
      <Footer />
    </div>
  );
}
