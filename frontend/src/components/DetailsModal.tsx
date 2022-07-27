import {
  Box,
  Button,
  Flex,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Tag,
  Text,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import React, { FC } from "react";
import { FiMenu, FiMoreVertical } from "react-icons/fi";
import { lunii } from "../../wailsjs/go/models";
import { DeleteModal } from "./DeleteModal";
import { PackTag } from "./PackTag";
import parse from "html-react-parser";
import { ListPacks, RemovePack } from "../../wailsjs/go/main/App";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useLocation, useRoute } from "wouter";

export const DetailsModal: FC = () => {
  const [, setLocation] = useLocation();
  const [, params] = useRoute("/pack/:uuid");
  const toast = useToast();
  const { data: packs } = useQuery(["packs"], ListPacks);

  const pack = packs?.find((p) => (p.uuid as any) === params?.uuid);
  if (!pack) return null;
  const queryClient = useQueryClient();

  const parsedDescription = parse(pack.description);

  const handleRemovePack = async () => {
    await RemovePack(pack.uuid);
    toast({
      title: "The pack was deleted from the device",
      status: "success",
      isClosable: true,
    });
    await queryClient.invalidateQueries(["device"]);
    await queryClient.invalidateQueries(["packs"]);
    setLocation("/");
  };

  return (
    <>
      <Modal isOpen={true} onClose={() => setLocation("/")}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader></ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Box mb={2}>
              <PackTag metadata={pack} />
            </Box>
            <Text fontSize={20} mb={2} fontWeight="bold">
              {pack.title}
            </Text>
            <Box mb={2}>{parsedDescription}</Box>
            <Tag mb={2}>{pack.uuid}</Tag>
          </ModalBody>

          <ModalFooter>
            <DeleteModal onDelete={handleRemovePack} />
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
