import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  Tooltip,
  useDisclosure,
} from "@chakra-ui/react";
import { useState } from "react";
import { MdCreate } from "react-icons/md";
import { CreatePack } from "../../wailsjs/go/main/App";

export const NewStoryModal = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [isLoading, setIsLoading] = useState(false);
  const handleCreatePack = async () => {
    setIsLoading(true);
    await CreatePack().catch((err) => console.log(err));
    setIsLoading(false);
  };
  return (
    <Button
      variant="outline"
      colorScheme="linkedin"
      rightIcon={<MdCreate />}
      ml={2}
      onClick={handleCreatePack}
      isLoading={isLoading}
    >
      Create story
    </Button>
  );
  return (
    <>
      <Tooltip label="Create a STUdio story pack from a directory">
        <Button
          variant="outline"
          colorScheme="linkedin"
          rightIcon={<MdCreate />}
          ml={2}
          onClick={onOpen}
        >
          Create story
        </Button>
      </Tooltip>

      <Modal isOpen={isOpen} onClose={onClose} isCentered size="2xl">
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Create a story pack</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>Choose a directory that follow the SimplePack structure</Text>
            <Text>Choose </Text>
          </ModalBody>

          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={onClose}>
              Close
            </Button>
            <Button variant="ghost">Secondary Action</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
