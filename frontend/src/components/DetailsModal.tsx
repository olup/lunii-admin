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
import {
  MdDetails,
  MdMenu,
  MdMore,
  MdOpenInFull,
  MdPlusOne,
  MdRemoveRedEye,
} from "react-icons/md";
import { lunii } from "../../wailsjs/go/models";
import { DeleteModal } from "./DeleteModal";
import { PackTag } from "./PackTag";
import parse from "html-react-parser";
import { ListPacks, RemovePack } from "../../wailsjs/go/main/App";
import { useQuery } from "@tanstack/react-query";

export const DetailsModal: FC<{
  uuid: number[];
}> = ({ uuid }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { data: packs, refetch: refetchPacks } = useQuery(["packs"], ListPacks);
  const pack = packs?.find((p) => p.uuid === uuid);
  const toast = useToast();

  if (!pack) return null;

  const parsedDescription = parse(pack.description);

  const handleRemovePack = async () => {
    await RemovePack(pack.uuid);
    toast({
      title: "The pack was deleted from the device",
      status: "success",
      duration: 9000,
      isClosable: true,
    });
    await refetchPacks();
  };

  return (
    <>
      <IconButton
        variant="ghost"
        aria-label="Details"
        icon={<MdMenu />}
        onClick={onOpen}
      />

      <Modal isOpen={isOpen} onClose={onClose}>
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
