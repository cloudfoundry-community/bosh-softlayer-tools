package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("task command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "task"}
		options = cmds.Options{
			Verbose: false,
			TaskID:  1,
			Debug:   false,
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewTaskCommand(options, fakeBmpClient)
	})

	Describe("NewTaskCommand", func() {
		It("create new TaskCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewTaskCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a TaskCommand", func() {
			Expect(cmd.Name()).To(Equal("task"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a TaskCommand", func() {
			Expect(cmd.Description()).To(Equal(`Show the output of the task: \"option --debug, Get the debug info of the task\"`))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a TaskCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp task <task-id>"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a TaskCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good TaskCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("validates a bad TaskCommand", func() {
			BeforeEach(func() {
				options = cmds.Options{
					Verbose: false,
					TaskID:  0,
					Debug:   false,
				}
			})

			It("no task_id", func() {
				cmd = bmp.NewLoginCommand(options, fakeBmpClient)

				validate, err := cmd.Validate()
				Expect(validate).To(BeFalse())
				Expect(err).To(HaveOccurred())
			})

		})
	})

	Describe("#Execute", func() {
		Context("executes a good TaskCommand", func() {
			Context("executes TaskCommand with default event level", func() {
				BeforeEach(func() {
					fakeBmpClient.TaskOutputResponse.Status = 200
					fakeBmpClient.TaskOutputErr = nil
					options = cmds.Options{
						Verbose: false,
						TaskID:  1,
						Debug:   false,
					}
				})

				It("executes a good TaskCommand ", func() {
					rc, err := cmd.Execute([]string{"bmp", "task", "--task_id=1"})
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("executes TaskCommand with debug level", func() {
				BeforeEach(func() {
					fakeBmpClient.TaskOutputResponse.Status = 200
					fakeBmpClient.TaskOutputErr = nil
					options = cmds.Options{
						Verbose: false,
						TaskID:  1,
						Debug:   true,
					}
				})

				It("executes a good TaskCommand ", func() {
					rc, err := cmd.Execute([]string{"bmp", "task", "--task_id=1", "--debug"})
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})

		})

		Context("executes a bad TaskCommand", func() {
			Context("executes TaskCommand with error", func() {
				BeforeEach(func() {
					fakeBmpClient.TaskOutputResponse.Status = 500
					fakeBmpClient.TaskOutputErr = errors.New("500")
				})

				It("executes with error", func() {
					rc, err := cmd.Execute([]string{"bmp", "task", "--task_id=1"})
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("TaskCommand response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.TaskOutputResponse.Status = 404
				})

				It("response code different than 200", func() {
					rc, err := cmd.Execute([]string{"bmp", "task", "--task_id=1"})
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
