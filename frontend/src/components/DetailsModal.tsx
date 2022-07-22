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
} from "@chakra-ui/react";
import React, { FC } from "react";
import { MdMenu } from "react-icons/md";
import { lunii } from "../../wailsjs/go/models";
import { DeleteModal } from "./DeleteModal";
import { PackTag } from "./PackTag";
import parse from "html-react-parser";

export const DetailsModal: FC<{
  metadata: lunii.Metadata;
  onDelete: () => any;
}> = ({ metadata, onDelete }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const parsedDescription = parse(metadata.description);

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
              <PackTag metadata={metadata} />
            </Box>
            <Text fontSize={20} mb={2} fontWeight="bold">
              {metadata.title}
            </Text>
            <Box mb={2}>{parsedDescription}</Box>
            <Tag mb={2}>{metadata.uuid}</Tag>
          </ModalBody>

          <ModalFooter>
            <DeleteModal onDelete={onDelete} />
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
