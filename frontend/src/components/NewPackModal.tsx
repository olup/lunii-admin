import {
  Alert,
  AlertIcon,
  Box,
  Button,
  Flex,
  Link,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useToast,
} from "@chakra-ui/react";
import { useState } from "react";
import { FiFile, FiFolder, FiPackage } from "react-icons/fi";
import { useLocation } from "wouter";
import { CreatePack, OpenDirectory, SaveFile } from "../../wailsjs/go/main/App";
import { BrowserOpenURL } from "../../wailsjs/runtime";
import { basename, dirname } from "../utils";

export const NewPackModal = () => {
  const toast = useToast();

  const [directoryPath, setDirectoryPath] = useState("");
  const [destinationPath, setDestinationPath] = useState("");
  const [_, setLocation] = useLocation();

  const handleClose = () => {
    setLocation("/");
  };

  return (
    <>
      <Modal isOpen={true} onClose={handleClose}>
        <ModalOverlay>
          <ModalContent>
            <ModalHeader fontSize="lg" fontWeight="bold">
              Creating a new pack
            </ModalHeader>

            <ModalBody>
              <Box>
                To create a new pack, first select a structured directory from
                your system
                <Alert variant="left-accent" colorScheme="gray" my={2}>
                  <AlertIcon />
                  <Text>
                    Not sure what is a structured directory ? Check the{" "}
                    <Link
                      textDecor="underline"
                      onClick={() =>
                        BrowserOpenURL(
                          "https://github.com/olup/lunii-admin#creating-your-own-story-packs"
                        )
                      }
                    >
                      project readme
                    </Link>
                  </Text>
                </Alert>
                <Flex mt={2}>
                  <Button
                    leftIcon={<FiFolder />}
                    colorScheme={directoryPath ? "green" : undefined}
                    onClick={() =>
                      OpenDirectory("Choose Directory").then(setDirectoryPath)
                    }
                  >
                    {directoryPath || "Choose directory"}
                  </Button>
                </Flex>
              </Box>
              {directoryPath && (
                <Box mt={2}>
                  Where should the pack be saved (.zip) ?
                  <Flex mt={2}>
                    <Button
                      leftIcon={<FiFile />}
                      colorScheme={destinationPath ? "green" : undefined}
                      onClick={() =>
                        SaveFile(
                          "Select destination",
                          dirname(directoryPath),
                          `${basename(directoryPath)}.zip`
                        ).then(setDestinationPath)
                      }
                    >
                      {destinationPath || "Select destination"}
                    </Button>
                  </Flex>
                </Box>
              )}
            </ModalBody>

            <ModalFooter>
              <Button onClick={handleClose}>Cancel</Button>
              {directoryPath && destinationPath && (
                <Button
                  rightIcon={<FiPackage />}
                  colorScheme="green"
                  ml={2}
                  onClick={() =>
                    CreatePack(directoryPath, destinationPath)
                      .then(() => {
                        toast({
                          title: "Your story pack was created",
                          description: "You can now install it on your device",
                          status: "success",
                          isClosable: true,
                        });
                        handleClose();
                      })
                      .catch((err) => {
                        toast({
                          title: "Something went wrong creating the pack",
                          description: err,
                          status: "error",
                          isClosable: true,
                        });
                        handleClose();
                      })
                  }
                >
                  Create
                </Button>
              )}
            </ModalFooter>
          </ModalContent>
        </ModalOverlay>
      </Modal>
    </>
  );
};
