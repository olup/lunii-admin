import { useTab, useToast } from "@chakra-ui/react";
import { useQueryClient, useMutation } from "@tanstack/react-query";
import { ChangePackOrder } from "../../wailsjs/go/main/App";
import { lunii } from "../../wailsjs/go/models";

const ChangePackOrderObject = ({
  id,
  position,
}: {
  id: string;
  position: number;
}) => ChangePackOrder(id, position);

export const useChangePackOrder = () => {
  const queryClient = useQueryClient();
  const toast = useToast();
  const mutation = useMutation(ChangePackOrderObject, {
    onMutate: async (arg) => {
      await queryClient.cancelQueries(["packs"]);
      const previousPacks = queryClient.getQueryData(["packs"]);

      queryClient.setQueryData(["packs"], (old) => {
        // optimistic reordering
        const newPacks = [...(old as lunii.Metadata[])];
        const index = newPacks.findIndex((p) => (p.uuid as any) === arg.id);
        var element = newPacks[index];
        newPacks.splice(index, 1);
        newPacks.splice(arg.position, 0, element);
        return newPacks;
      });
      return { previousPacks };
    },
    onError: (err, newTodo, context) => {
      queryClient.setQueryData(["packs"], (context as any).previousPacks);
      toast({
        status: "error",
        title: "Could not reorder the pack",
        description: err as string,
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries(["packs"]);
    },
  });

  return mutation;
};
