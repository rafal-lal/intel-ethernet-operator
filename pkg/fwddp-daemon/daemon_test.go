// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package daemon

import (
	"context"
	"os/exec"
	"path"
	"syscall"

	gerrors "errors"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ethernetv1 "github.com/otcshare/intel-ethernet-operator/apis/ethernet/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type TestData struct {
	NodeConfig ethernetv1.EthernetNodeConfig
	Inventory  []ethernetv1.Device
	Node       core.Node
}

func (d *TestData) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: d.NodeConfig.Namespace,
		Name:      d.NodeConfig.Name,
	}
}

func initTestData() TestData {
	return TestData{
		NodeConfig: ethernetv1.EthernetNodeConfig{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: ethernetv1.EthernetNodeConfigSpec{
				Config: []ethernetv1.DeviceNodeConfig{
					{
						PCIAddress: "00:00:00.1",
						DeviceConfig: ethernetv1.DeviceConfig{
							DDPURL: "http://testddpurl",
							FWURL:  "http://testfwurl",
						},
					},
				},
			},
		},
		Node: core.Node{
			ObjectMeta: v1.ObjectMeta{
				Name:   "test",
				Labels: map[string]string{"fpga.ethernet.com/intel-ethernet-present": ""},
			},
		},
		Inventory: []ethernetv1.Device{
			{
				PCIAddress:    "00:00:00.0",
				Name:          "TestName",
				Driver:        "TestDriver",
				DriverVersion: "TestDriverVersion",
				Firmware: ethernetv1.FirmwareInfo{
					MAC:     "aa:bb:cc:dd:ee:ff",
					Version: "TestFWVersion",
				},
			},
		},
	}
}

func initReconciler(toBeInitialized *NodeConfigReconciler, nodeName, namespace string) error {
	cset, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	r, err := NewNodeConfigReconciler(k8sClient, cset, log, nodeName, namespace)
	if err != nil {
		return err
	}

	*toBeInitialized = *r
	return nil
}

var data = TestData{}

var _ = Describe("FirmwareDaemonTest", func() {
	reconciler := new(NodeConfigReconciler)
	var _ = BeforeEach(func() {
		data = initTestData()
		compatMapPath = "testdata/supported_devices.json"

		getInventory = func(_ logr.Logger) ([]ethernetv1.Device, error) {
			return data.Inventory, nil
		}
		utilsDownloadFile = func(path, url, checksum string, _ logr.Logger) error {
			return nil
		}
		utilsUntar = func(srcPath string, dstPath string, log logr.Logger) error {
			return nil
		}
		nvmupdateExec = func(cmd *exec.Cmd, log logr.Logger) error {
			return nil
		}
		fwInstallDest = "./workdir/nvmupdate/"
	})

	var _ = Context("Reconciler", func() {
		BeforeEach(func() {
		})

		AfterEach(func() {
			nn := data.GetNamespacedName()
			if err := k8sClient.Get(context.TODO(), nn, &data.NodeConfig); err == nil {
				data.NodeConfig.Spec = ethernetv1.EthernetNodeConfigSpec{
					Config: []ethernetv1.DeviceNodeConfig{},
				}
				Expect(k8sClient.Update(context.TODO(), &data.NodeConfig)).NotTo(HaveOccurred())
				_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nn})
				Expect(err).ToNot(HaveOccurred())
				Expect(k8sClient.Delete(context.TODO(), &data.NodeConfig)).ToNot(HaveOccurred())
			} else if errors.IsNotFound(err) {
				log.Info("Requested NodeConfig does not exists", "NodeConfig", &data.NodeConfig)
			} else {
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(k8sClient.Delete(context.TODO(), &data.Node)).To(Succeed())
		})

		var _ = It("will create empty NodeConfig if not exits", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: data.GetNamespacedName()})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Devices).To(HaveLen(0))
		})

		var _ = It("will update inventory on Reconcile()", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: data.GetNamespacedName()})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Devices).To(HaveLen(0))

			_, err = reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: data.GetNamespacedName()})
			Expect(err).ToNot(HaveOccurred())

			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Devices).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Devices[0]).To(Equal(data.Inventory[0]))

		})

		var _ = It("will ignore CRs with wrong name", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      "othername",
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(0))
		})

		var _ = It("will ignore CRs with wrong namespace", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "othernamespace",
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(0))
		})

		var _ = It("will update condition to Inventory up to date if Spec.Config is empty", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())

			data.NodeConfig.Spec.Config = []ethernetv1.DeviceNodeConfig{}

			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())
			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateNotRequested)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal("Inventory up to date"))
		})

		var _ = It("will update condition to UpdateFailed if no matching devices were found", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())
			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())
			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateFailed)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal("Device 00:00:00.1 not found"))
		})

		var _ = It("will update condition to UpdateFailed if not able to download firmware", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())
			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())

			data.Inventory[0].PCIAddress = "00:00:00.1"

			downloadErr := gerrors.New("Unable to download")
			utilsDownloadFile = func(path, url, checksum string, _ logr.Logger) error {
				return downloadErr
			}

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateFailed)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal(downloadErr.Error()))
		})

		var _ = It("will update condition to UpdateFailed if not able to untar firmware", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())
			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())

			data.Inventory[0].PCIAddress = "00:00:00.1"

			untarErr := gerrors.New("Unable to untar")
			utilsUntar = func(srcPath string, dstPath string, log logr.Logger) error {
				return untarErr
			}

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateFailed)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal(untarErr.Error()))
		})

		var _ = It("will update condition to UpdateFailed if firmware update fails", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())
			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())

			data.Inventory[0].PCIAddress = "00:00:00.1"

			fwErr := gerrors.New("Unable to update firmware")
			nvmupdateExec = func(cmd *exec.Cmd, log logr.Logger) error {
				return fwErr
			}

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateFailed)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal(fwErr.Error()))
		})

		var _ = It("will update update condition to UpdateSucceeded after successful firmware update", func() {
			Expect(k8sClient.Create(context.TODO(), &data.Node)).To(Succeed())
			Expect(k8sClient.Create(context.TODO(), &data.NodeConfig)).To(Succeed())

			data.Inventory[0].PCIAddress = "00:00:00.1"

			rootAttr := &syscall.SysProcAttr{
				Credential: &syscall.Credential{Uid: 0, Gid: 0},
			}
			nvmupdateExec = func(cmd *exec.Cmd, log logr.Logger) error {
				Expect(cmd.SysProcAttr).To(Equal(rootAttr))
				Expect(cmd.Dir).To(Equal(path.Join(fwInstallDest, data.NodeConfig.Spec.Config[0].PCIAddress,
					nvmupdate64eDirSuffix)))
				return nil
			}

			Expect(initReconciler(reconciler, data.NodeConfig.Name, data.NodeConfig.Namespace)).To(Succeed())

			_, err := reconciler.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: data.NodeConfig.Namespace,
				Name:      data.NodeConfig.Name,
			}})
			Expect(err).ToNot(HaveOccurred())

			nodeConfigs := &ethernetv1.EthernetNodeConfigList{}
			Expect(k8sClient.List(context.TODO(), nodeConfigs)).To(Succeed())
			Expect(nodeConfigs.Items).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions).To(HaveLen(1))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Reason).To(Equal(string(UpdateSucceeded)))
			Expect(nodeConfigs.Items[0].Status.Conditions[0].Message).To(Equal("Updated successfully"))
		})
	})
})
