"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { userService } from "@/libs/services";
import { getSessionUser, clearSession } from "@/libs/session";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";
import type { User } from "@/types/services/user";

export function DeleteAccount() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [confirmation, setConfirmation] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    async function loadUser() {
      const session = await getSessionUser();
      if (!session) {
        router.push("/user/register");
        return;
      }
      setUser(session);
    }
    loadUser();
  }, [router]);

  async function handleDelete(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (confirmation !== "EXCLUIR") {
      setError("Digite EXCLUIR para confirmar");
      return;
    }

    if (!user) return;

    setLoading(true);
    const success = await userService.delete(user.id);
    if (success) {
      clearSession();
      router.push("/");
    } else {
      setError("Erro ao excluir conta. Tente novamente.");
    }
    setLoading(false);
  }

  if (!user) return null;

  return (
    <FormPanel title="Excluir conta" maxWidth="max-w-md">
      <p className="text-sm text-neutral-slate mb-4">
        Esta ação é <strong className="text-red-400">irreversível</strong>.
        Todos os seus dados serão permanentemente removidos.
      </p>

      <Alert variant="error">
        Você está prestes a excluir a conta de{" "}
        <strong className="text-white">{user.email}</strong>.
      </Alert>

      {error && <Alert variant="error">{error}</Alert>}

      <form onSubmit={handleDelete} className="flex flex-col gap-5">
        <FormInput
          label={`Digite EXCLUIR para confirmar`}
          name="confirmation"
          value={confirmation}
          onChange={(e) => {
            setConfirmation(e.target.value);
            setError("");
          }}
          placeholder="EXCLUIR"
        />

        <div className="flex gap-4">
          <FormButton
            type="button"
            variant="outline"
            onClick={() => router.back()}
            fullWidth
          >
            Cancelar
          </FormButton>
          <FormButton
            type="submit"
            variant="danger"
            loading={loading}
            loadingText="Excluindo..."
            disabled={confirmation !== "EXCLUIR"}
            fullWidth
          >
            Excluir conta
          </FormButton>
        </div>
      </form>
    </FormPanel>
  );
}
