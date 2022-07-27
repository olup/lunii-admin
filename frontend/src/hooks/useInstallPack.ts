import { useToast } from "@chakra-ui/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { OpenFile, InstallPack, ListPacks } from "../../wailsjs/go/main/App";

export const useInstallPack = () => {
  const toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const queryClient = useQueryClient();

  const mutate = async () => {
    const packPath = await OpenFile("Select your pack");
    if (!packPath) return;

    setIsLoading(true);
    try {
      await InstallPack(packPath);
      toast({
        title: "The pack was installed on the device",
        status: "success",
        isClosable: true,
      });
      await queryClient.invalidateQueries(["device"]);
      await queryClient.invalidateQueries(["packs"]);
    } catch (e) {
      toast({
        title: "Could not install pack on device",
        status: "warning",
        description:
          "Something went wrong while installing the pack on your device",
      });
    }
    setIsLoading(false);
  };
  return { mutate, isLoading };
};
