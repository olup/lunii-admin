import { useToast } from "@chakra-ui/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { OpenFiles, InstallPack, ListPacks } from "../../wailsjs/go/main/App";

export const useInstallPack = () => {
  const toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const queryClient = useQueryClient();

  const mutate = async () => {
    const packPaths = await OpenFiles("Select your pack");
    if (!packPaths?.length) return;

    toast({
      title: `Installing ${packPaths.length} pack(s) on the device`,
      status: "info",
      isClosable: true,
      duration: 1000,
    });

    setIsLoading(true);
    for (const packPath of packPaths) {
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
    }
    setIsLoading(false);
  };
  return { mutate, isLoading };
};
